package algo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

const DEFAULT_PREFIX = "no_prefix"
const SEP = "_"
const MIN_ID_LEN = 1

var STRICT_MODE = true

/* To16Bytes expects string format like {entity}_{prefix}_{raw}
 * entity is ignored when calculating hash
 * prefix hash to 4 bytes
 * raw split to 3 parts and hash to 12 bytes
 * DEFAULT_PREFIX used when SEP not found in source string
 */
func To16Bytes(s string) (bytes []byte, err error) {
	if len(s) < MIN_ID_LEN {
		err = errors.New(fmt.Sprintf("Input string is at least length of %s", MIN_ID_LEN))
		return
	}
	i := strings.Index(s, SEP)
	if i > -1 {
		s = s[i+1:]
		// log.Println("Ignoring entity... The left is " + s)
	} else if STRICT_MODE {
		err = errors.New(fmt.Sprintf("No separator '%s' found in string param", SEP))
		return
	}

	i = strings.Index(s, SEP)
	var prefix, raw string

	if i > -1 {
		prefix = s[:i]
		raw = s[i+1:]
	} else if STRICT_MODE {
		err = errors.New(fmt.Sprintf("Should have two '%s' in string param", SEP))
		return
	} else {
		prefix = DEFAULT_PREFIX
		raw = s
	}

	if len(raw) < 1 {
		err = errors.New("Length of raw id is 0.")
		return
	}

	var pureNum bool = true
	var ui64 uint64
	if ui64, err = strconv.ParseUint(raw, 10, 64); err != nil {
		pureNum = false
	}

	bytes = make([]byte, 16)
	h := fnv.New32a() // general fast hash

	if pureNum { // pure number like github

		binary.LittleEndian.PutUint64(bytes[:8], ui64)

		/* } else if raw is uuid {
		 * uuid specific hash
		 * }
		 */

	} else { // other raw string
		// divide raw into 3 parts, hash to 4 bytes each part
		var last int
		n := len(raw) / 3
		for i := 0; i < 3; i++ {
			h.Reset()
			if i == 2 {
				last = len(raw)
			} else {
				last = n * (i + 1)
			}
			h.Write([]byte(raw[n*i : last]))
			// put less changed behind
			binary.LittleEndian.PutUint32(bytes[(3-1-i)*4:], h.Sum32())
		}

	}

	// last 4 bytes for prefix, put less changed behind
	h.Reset()
	h.Write([]byte(prefix))
	binary.LittleEndian.PutUint32(bytes[12:], h.Sum32())

	return
}
