package internal

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	parser "golang.org/x/net/html"
)

const downloadPage = "https://sqlite.org/download.html"

func MustParse() []SqliteProduct {
	products, _ := getProductCsv(downloadPage)

	return products
}

func getProductCsv(link string) ([]SqliteProduct, error) {
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	// content, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// html := string(content)

	t := parser.NewTokenizer(res.Body)
	sqliteProducts := make([]SqliteProduct, 0)
	for {
		// Moves the cursor to the next node in the tree.
		tt := t.Next()

		// html.ErrorToken is an error token.
		if tt == parser.ErrorToken {
			break
		}

		// Returns a token for the current token.
		tc := t.Token()

		// If the token is a CommentToken.
		if tc.Type == parser.CommentToken {
			if !strings.Contains(tc.Data, "Download product data for scripts to read") {
				continue
			}

			csvString := strings.Replace(tc.Data, "Download product data for scripts to read", "", 1)
			csvString = strings.TrimSpace(csvString)
			csvStringReader := strings.NewReader(csvString)

			products := csvToMap(csvStringReader)

			productNameRegexp := regexp.MustCompile(`sqlite\-([\-\w]+)\-\d+`)
			extensionRegexp := regexp.MustCompile(`\.([\.\w]+)$`)

			for _, product := range products {

				relativeUrl := product["RELATIVE-URL"]
				nameMatches := productNameRegexp.FindStringSubmatch(relativeUrl)
				extensionMatches := extensionRegexp.FindStringSubmatch(relativeUrl)

				sizeInBytes, err := strconv.ParseInt(product["SIZE-IN-BYTES"], 10, 64)
				if err != nil {
					return nil, err
				}

				sqlProduct := SqliteProduct{
					Name:        nameMatches[1],
					Extension:   extensionMatches[1],
					Version:     product["VERSION"],
					RelativeUrl: product["RELATIVE-URL"],
					SizeInBytes: sizeInBytes,
					Sha3Hash:    product["SHA3-HASH"],
				}

				sqliteProducts = append(sqliteProducts, sqlProduct)
			}
		}

	}

	return sqliteProducts, nil
}

// PRODUCT,VERSION,RELATIVE-URL,SIZE-IN-BYTES,SHA3-HASH
// PRODUCT,2023-05-02 16:34 UTC,snapshot/sqlite-snapshot-202305021634.tar.gz,3136585,73bcd248a08fe8937f6b9c838fd8c8badbf9a3e4954e331b0ef6287768e7352d
// PRODUCT,3.41.2,2023/sqlite-amalgamation-3410200.zip,2623078,c51ca72411b8453c64e0980be23bc9b9530bdc3ec1513e06fbf022ed0fd02463
// PRODUCT,3.41.2,2023/sqlite-autoconf-3410200.tar.gz,3125545,1ebb5539dd6fde9a0f89e8ab765af0b9f02010fc6baf6985b54781a38c00020a
// PRODUCT,3.41.2,2023/sqlite-doc-3410200.zip,10633275,8b3daa86d41cbc407e0992e11193d81094f01771c2dbeae93695f3dbaf19cbfe
// PRODUCT,3.41.2,2023/sqlite-android-3410200.aar,3420407,6cd0ca943868978a945d4531bdc4e21fee8e3a7650e577f052550a37bc62db6e
// PRODUCT,3.41.2,2023/sqlite-tools-linux-x86-3410200.zip,2269374,c41d2d97b62af6c6bd523ce9a06d1ea9f4e22e649f7a677a1c43f36d5ac3272c
// PRODUCT,3.41.2,2023/sqlite-tools-osx-x86-3410200.zip,1616529,0e0f1eeade6ff7bc6f0ddbf0fe86798a595494328e1afa15da192e5fc516cb29
// PRODUCT,3.41.2,2023/sqlite-dll-win32-x86-3410200.zip,576595,2a3b03595ee7c3b172465d9be3990fd8ac10fa93b9afd22031b67fa42c2325e1
// PRODUCT,3.41.2,2023/sqlite-dll-win64-x64-3410200.zip,1202571,a85aa4c03f394147b4a9050245fb8d5d0b00ae726d1f259e71c07055e44854cf
// PRODUCT,3.41.2,2023/sqlite-tools-win32-x86-3410200.zip,1999045,0ceebb7f8378707d6d6b0737ecdf2ba02253a3b44b1009400f86273719d98f1f
// PRODUCT,3.41.2,2023/sqlite-wasm-3410200.zip,719097,1713c3a4c25b09f5d52ed25fb81948f6081e9eeefb5a6d449aead1cd848e962c
// PRODUCT,3.41.2,2023/sqlite-src-3410200.zip,13836230,793e24c5158bafdb8a6e861ea9ad267cc773705d0bebb78b5c4aa323033e8028
// PRODUCT,3.41.2,2023/sqlite-preprocessed-3410200.zip,2731397,6ffb450572964d829c688c3767448eac89f3577adfb1f08bfd3aaf0f050c3596

func csvToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}
