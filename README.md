# IpInfo - selfhosted ip information service

Service can:
- Collect and update ip address information from RIR

Service provide:
- Endpoint for get information by ip

## Run service
```
make init
make build
make run
```

## Using
```
// Ip v4
GET /ipv4/127.0.0.1

// Ip v6
GET /ipv6/::1
```
## How it works

1. Get info from all 5 top-level RIR(Regional Internet Registries)
    1. ARIN(American Registry for Internet Numbers)
    2. RIPE NCC(Reseaux IP Europeens)
    3. APNIC(Asia-Pacific Network Information Centre)
    4. LACNIC(Latin America and Caribbean Network Information Centre)
    5. AFRINIC(African Network Information Centre)
2. Merge with previous data in own database

P.S.
Map of RIRs

![Map](./docs/rir-map.svg)

Url rules
```
// registry - one of { apnic, arin, iana, lacnic, ripencc }, domain changing too
// For more info visit https://www.iana.org/numbers
http://www.apnic.net/stats/<registry>/delegated-<registry>-latest
rsync www.apnic.net::/stats/<registry>/delegated-<registry>-latest
ftp://ftp.apnic.net/pub/stats/<registry>/delegated-<registry>-latest
```
Be careful. Slow downloading! Possible be more than 10 minutes.

In the same directory, you can find README file with specification of data format and more useful info.

Some info from README

1. The Version line:

    Format:
        `version|registry|serial|records|startdate|enddate|UTCoffset`
    - __version__ = format version number of this file, currently 2;
    - __registry__ = as for records and filename (see below);
    - __serial__ = serial number of this file (within the creating RIR series);
    - __records__ = number of records in file, excluding blank lines,
                     summary lines, the version line and comments;
    - __startdate__ = start date of time period, in yyyymmdd format;
    - __enddate__ = end date of period, in yyyymmdd format;
    - __UTCoffset__ = offset from UTC of local RIR producing file, 
                     in +/- HHMM format.

2. The Summary line:

    The summary lines count the number of record lines of each type in the file.

    Format:
        `registry|*|type|*|count|summary`
    - __registry__ = as for records (see below);
    - __\*__ = an ASCII '*' (unused field, retained for spreadsheet purposes);
    - __type__ = as for records (defined below);
    - __count__ = the number of record lines of this type in the file;
    - __summary__ = the ASCII string 'summary' (to distinguish the record line).


    Note that the count does not equate to the total amount of
    resources for each class of record. This is to be computed from
    the records themselves.

3. Record format:

    After the defined file header, and excluding any space or
    comments, each line in the file represents a single allocation
    (or assignment) of a specific range of Internet number resources
    (IPv4, IPv6 or ASN), made by the RIR identified in the record.

    In the case of IPv4 the records may represent non-CIDR ranges
    or CIDR blocks, and therefore the record format represents a
    beginning of range, and a count. This can be converted to
    prefix/length using simple algorithms.

    In the case of IPv6 the record format represents the prefix
    and the count of /128 instances under that prefix.

    Format:
        `registry|cc|type|start|value|date|status[|extensions...]`

    - __registry__ = One value from the set of defined strings: `{apnic,arin,iana,lacnic,ripencc}`.

    - __cc__ = ISO 2-letter country code of the organization to which the allocation or assignment was made, and the enumerated variances of `{AP,EU,UK}`. These values are not defined in ISO 3166 but are widely used.

    - __type__ = Type of Internet number resource represented in this record, One value from the set of defined strings: `{asn,ipv4,ipv6}`.

    - __start__ = In the case of records of type 'ipv4' or 'ipv6' this is the IPv4 or IPv6 'first' address' of the range.

		In the case of an 16 bit AS number  the format is
		the integer value in the range 0 to 65535, in the
		case of a 32 bit ASN the value is in the range 0
		to 4294967296. No distinction is drawn between 16
		and 32 bit ASN values in the range 0 to 65535.

    - __value__ = In the case of IPv4 address the count of hosts for
                this range. This count does not have to represent a
                CIDR range.

        In the case of an IPv6 address the value will be
        the CIDR prefix length from the 'first address'
        value of <start>.
        In the case of records of type 'asn' the number is
        the count of AS from this start value.

    - __date__ = Date on this allocation/assignment was made by the
                RIR in the format YYYYMMDD;

        Where the allocation or assignment has been
        transferred from another registry, this date
        represents the date of first assignment or allocation
        as received in from the original RIR.

        It is noted that where records do not show a date of
        first assignment, this can take the 00000000 value.

    - __status__ = Type of allocation from the set: `{allocated, assigned}`

        This is the allocation or assignment made by the
        registry producing the file and not any sub-assignment
        by other agencies. 

    - __extensions__ = Any extra data on a line is undefined, but should conform to use of the field separator used above.

## Preferences

Databases:
- ClickHouse
- PostgreSQL

Cache:
- Valkey
- Redis