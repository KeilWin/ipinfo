package app

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/dto/database"
)

type Rir struct {
	Domain   string
	PathName string
	FileName string
	DbName   string
}

func NewRir(domain, pathName, fileName, dbName string) *Rir {
	return &Rir{
		Domain:   domain,
		PathName: pathName,
		FileName: fileName,
		DbName:   dbName,
	}
}

var Rirs = [5]*Rir{
	NewRir("arin", "arin", "arin-extended", "arin"),
	NewRir("apnic", "apnic", "apnic", "apnic"),
	NewRir("afrinic", "afrinic", "afrinic", "afrinic"),
	NewRir("lacnic", "lacnic", "lacnic", "lacnic"),
	NewRir("ripe", "ripencc", "ripencc", "ripencc"),
}

func FindRirByDbName(name string) int {
	for i, rir := range Rirs {
		if rir.DbName == name {
			return i
		}
	}

	return -1
}

var IpVersions = [2]string{
	"ipv4",
	"ipv6",
}

var Statuses = [4]string{
	"allocated",
	"assigned",
	"available",
	"reserved",
}

const UnknownStatus = 5

func newDownloadUrl(domain, pathName, fileName string) string {
	return fmt.Sprintf("https://ftp.%s.net/pub/stats/%s/delegated-%s-latest", domain, pathName, fileName)
}

type RirManager struct {
	rir          *Rir
	downloadUrl  string
	downloadPath string
}

func (p *RirManager) Update(db database.Database, tableName string) error {
	data, err := p.Download()
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}

	err = p.Upload(db, tableName, data)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func NewEndRangeIpAddressV4(addrStart netip.Addr, quantity uint32) *netip.Addr {
	buf := addrStart.As4()
	addrNumber := uint32(0)
	for i := 0; i < 4; i++ {
		addrNumber |= uint32(buf[i]) << (32 - 8*(i+1))
	}
	addrNumber += quantity
	nBuff := [4]byte{}
	for i := 0; i < 4; i++ {
		nBuff[i] = byte(addrNumber >> (32 - 8*(i+1)))
	}
	endAddr := netip.AddrFrom4(nBuff)
	return &endAddr
}

func NewEndRangeIpAddressV6(addrStart netip.Addr, quantity uint64) (*netip.Addr, error) {
	addrNumber := new(big.Int)
	addrNumber.SetBytes(addrStart.AsSlice())
	addrNumber.Add(addrNumber, big.NewInt(int64(quantity)))
	bs := addrNumber.Bytes()
	buf := make([]byte, 16)
	copy(buf[len(buf)-len(bs):], bs)
	addrEnd, ok := netip.AddrFromSlice(buf)
	if !ok {
		return nil, fmt.Errorf("can't convert buffer to address ipv6: %v", buf)
	}

	return &addrEnd, nil
}

func (p *RirManager) Download() ([]common.IpRange, error) {
	slog.Info("starting download", "rir", p.rir.DbName, "url", p.downloadUrl)
	cli := http.Client{
		Timeout: 600 * time.Second,
	}
	response, err := cli.Get(p.downloadUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", response.Status)
	}

	reader := bufio.NewReaderSize(response.Body, 1<<20)

	var line string
	for {
		line, err = reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		break
	}

	for i := 0; i < 3; i++ {
		line, err = reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
	}

	slog.Info("start readind")
	ipRanges := make([]common.IpRange, 0, 5000)

	for {
		line, err = reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		valArray := strings.Split(line, "|")

		rirId := FindRirByDbName(valArray[0])
		if rirId == -1 {
			return nil, fmt.Errorf("can't parse rirId from line: %s", line)
		}

		if valArray[2] == "asn" {
			continue
		}
		versionIpId := slices.Index(IpVersions[:], valArray[2])
		if versionIpId == -1 {
			return nil, fmt.Errorf("can't parse versionIpId = '%s' from line: %s", valArray[2], line)
		}

		// better uint64 or big.Int
		quantity, err := strconv.ParseUint(valArray[4], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("can't parse quantity: %w", err)
		}

		addrStart, err := netip.ParseAddr(valArray[3])
		if err != nil {
			return nil, fmt.Errorf("can't parse ip: %w", err)
		}
		var addrEnd *netip.Addr
		if addrStart.Is4() {
			addrEnd = NewEndRangeIpAddressV4(addrStart, uint32(quantity))
		} else if addrStart.Is6() {
			addrEnd, err = NewEndRangeIpAddressV6(addrStart, quantity)
			if err != nil {
				return nil, fmt.Errorf("can't compute addrEnd: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unknown ip format = '%s' from line: %s", addrStart.String(), line)
		}

		var createdAt sql.NullTime
		if valArray[5] == "" {
			createdAt = sql.NullTime{Valid: false}
		} else {
			date, err := time.Parse("20060102", valArray[5])
			if err != nil {
				return nil, fmt.Errorf("can't parse date = '%s' from line: %s", valArray[5], line)
			}
			createdAt = sql.NullTime{Time: date, Valid: true}
		}

		statusId := slices.Index(Statuses[:], valArray[6])
		if statusId == -1 {
			slog.Info("unknown status = '%s' from line: %s", valArray[6], line)
			statusId = UnknownStatus
		}

		ipRanges = append(ipRanges, common.IpRange{
			Rir:         rirId + 1,
			CountryCode: valArray[1],
			VersionIp:   versionIpId + 1,
			StartIp:     valArray[3],
			EndIp:       addrEnd.String(),
			Quantity:    quantity,
			Status:      statusId + 1,
			CreatedAt:   createdAt,
		})
	}

	return ipRanges, nil
}

func (p *RirManager) Upload(db database.Database, tableName string, data []common.IpRange) error {
	return db.CopyToIpRangesFromArray(tableName, data)
}

func NewRirManager(rir *Rir) *RirManager {
	return &RirManager{
		rir:         rir,
		downloadUrl: newDownloadUrl(rir.Domain, rir.PathName, rir.FileName),
	}
}
