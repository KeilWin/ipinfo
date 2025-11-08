package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/KeilWin/ipinfo/internal/common"
	"github.com/KeilWin/ipinfo/internal/dto/database"
)

var RirsNames = [5]string{
	"apnic",
	// "arin",
	// "iana",
	// "lacnic",
	// "ripencc",
}
var IpVersions = [2]string{
	"ipv4",
	"ipv6",
}

var Statuses = [2]string{
	"allocated",
	"assigned",
}

type Rir struct {
	name         string
	downloadUrl  string
	downloadPath string
}

func (p *Rir) Update(db database.Database, tableName string) error {
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

func (p *Rir) Download() ([]common.IpRange, error) {
	slog.Info("starting download", "rir", p.name, "url", p.downloadUrl)
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

		rirId := slices.Index(RirsNames[:], valArray[0])
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

		_, err := strconv.Atoi(valArray[4])
		if err != nil {
			return nil, fmt.Errorf("can't parse quantity: %w", err)
		}

		statusId := slices.Index(Statuses[:], valArray[6])
		if statusId == -1 {
			return nil, fmt.Errorf("can't parse statusId = '%s' from line: %s", valArray[6], line)
		}

		ipRanges = append(ipRanges, common.IpRange{
			Rir:         rirId + 1,
			CountryCode: valArray[1],
			VersionIp:   versionIpId + 1,
			StartIp:     valArray[3],
			EndIp:       valArray[3],
			Status:      statusId + 1,
			CreatedAt:   valArray[5],
		})
	}

	return ipRanges, nil
}

func (p *Rir) Upload(db database.Database, tableName string, data []common.IpRange) error {
	return db.CopyToIpRangesFromArray(tableName, data)
}

func NewRir(name, path string) *Rir {
	return &Rir{
		name:         name,
		downloadUrl:  newDownloadUrl(name),
		downloadPath: newDownloadPath(name, path),
	}
}

func newDownloadUrl(name string) string {
	return fmt.Sprintf("https://ftp.%s.net/pub/stats/%s/delegated-%s-latest", name, name, name)
}

func newDownloadPath(name, path string) string {
	return fmt.Sprintf("%s/%s", path, name)
}
