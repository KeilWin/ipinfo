package app

import (
	"bufio"
	"context"
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

func newDownloadUrl(rir *Rir) string {
	return fmt.Sprintf("https://ftp.%s.net/pub/stats/%s/delegated-%s-latest", rir.Domain, rir.PathName, rir.FileName)
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

type RirManager struct {
	Rir          *Rir
	db           database.Database
	ctx          context.Context
	timeToUpdate time.Time
}

func (p *RirManager) GetLastUpdate() (*time.Time, error) {
	var lastUpdateDateTime time.Time
	lastUpdateOption, err := p.db.GetOption(fmt.Sprintf("lastUpdate%s", p.Rir.DbName), p.ctx)
	if errors.Is(err, sql.ErrNoRows) {
		lastUpdateDateTime = time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	} else if err != nil {
		return nil, fmt.Errorf("get option 'lastUpdate': %w", err)
	} else {
		lastUpdateDateTime, err = time.ParseInLocation("2006-01-02 15:04:05", lastUpdateOption, time.UTC)
		if err != nil {
			return nil, fmt.Errorf("parse option 'lastUpdate': %w", err)
		}
	}
	return &lastUpdateDateTime, nil
}

func (p *RirManager) RefreshLastUpdate() (time.Time, error) {
	now := time.Now().UTC()
	return now, p.db.UpdateOption(fmt.Sprintf("lastUpdate%s", p.Rir.DbName), now.Format("2006-01-02 15:04:05"), p.ctx)

}

func (p *RirManager) Start() error {
	slog.Info("rir manager started", "rir", p.Rir.DbName)
	lastUpdate, err := p.GetLastUpdate()
	if err != nil {
		return fmt.Errorf("lastUpdate: %w", err)
	}
	slog.Info("get lastUpdate ok", "lastUpdate", lastUpdate.Format("2006-01-02 15:04:05"))

	now := time.Now().UTC()
	if now.Sub(*lastUpdate).Hours() >= 24 {
		slog.Info("updating", "rir", p.Rir.DbName)
		data, err := p.Download()
		if err != nil {
			return fmt.Errorf("download: %w", err)
		}

		rows, err := p.ParseData(data)
		if err != nil {
			return fmt.Errorf("parse data: %w", err)
		}

		err = p.Upload(rows)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}

		nowDt, err := p.RefreshLastUpdate()
		if err != nil {
			return fmt.Errorf("refresh lastUpdate: %w", err)
		}

		slog.Info("successful update", "rir", p.Rir.DbName)
		time.Sleep(time.Until(time.Date(nowDt.Year(), nowDt.Month(), nowDt.Day(), p.timeToUpdate.Hour(), p.timeToUpdate.Minute(), p.timeToUpdate.Second()+5, p.timeToUpdate.Nanosecond(), time.UTC).Add(24 * time.Hour)))
	} else {
		slog.Info("wake up too early", "rir", p.Rir.DbName)
		time.Sleep(now.Sub(*lastUpdate) + 5*time.Second)
	}

	return nil
}

func (p *RirManager) ParseHeader(reader *bufio.Reader) error {
	var err error
	var line string

	for {
		line, err = reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
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
			return err
		}
	}
	return nil
}

func (p *RirManager) ParseData(data io.ReadCloser) ([]common.IpRange, error) {
	var err error
	var line string

	reader := bufio.NewReaderSize(data, 1<<20)
	if err = p.ParseHeader(reader); err != nil {
		return nil, fmt.Errorf("parse header: %w", err)
	}

	ipRanges := make([]common.IpRange, 0, 10000)
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

		if valArray[2] == "asn" {
			continue
		}

		rirId := FindRirByDbName(valArray[0])
		if rirId == -1 {
			return nil, fmt.Errorf("can't parse rirId from line: %s", line)
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

		var statusChangedAt sql.NullTime
		if valArray[5] == "" {
			statusChangedAt = sql.NullTime{Valid: false}
		} else {
			date, err := time.Parse("20060102", valArray[5])
			if err != nil {
				return nil, fmt.Errorf("can't parse date = '%s' from line: %s", valArray[5], line)
			}
			statusChangedAt = sql.NullTime{Time: date, Valid: true}
		}

		statusId := slices.Index(Statuses[:], valArray[6])
		if statusId == -1 {
			slog.Info("unknown status = '%s' from line: %s", valArray[6], line)
			statusId = UnknownStatus
		}

		ipRanges = append(ipRanges, common.IpRange{
			RirId:           rirId + 1,
			CountryCode:     valArray[1],
			IpVersionId:     versionIpId + 1,
			StartIp:         valArray[3],
			EndIp:           addrEnd.String(),
			Quantity:        quantity,
			StatusId:        statusId + 1,
			StatusChangedAt: statusChangedAt,
		})
	}

	return ipRanges, nil
}

func (p *RirManager) Download() (io.ReadCloser, error) {
	cli := http.Client{
		Timeout: 600 * time.Second,
	}
	response, err := cli.Get(newDownloadUrl(p.Rir))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %s", response.Status)
	}

	return response.Body, nil
}

func (p *RirManager) Upload(data []common.IpRange) error {
	return p.db.UpdateRirData(p.Rir.DbName, data, p.ctx)
}

func NewRirManager(rir *Rir, db database.Database, ctx context.Context, timeToUpdate time.Time) *RirManager {
	return &RirManager{
		Rir:          rir,
		db:           db,
		ctx:          ctx,
		timeToUpdate: timeToUpdate,
	}
}
