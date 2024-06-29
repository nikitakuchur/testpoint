package writer

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"restcompare/internal/sender"
)

func WriteResponses(input <-chan sender.Response, dir string) {
	fileMap := make(map[string]*os.File)

	defer func() {
		for _, file := range fileMap {
			err := file.Sync()
			if err != nil {
				log.Fatalln("cannot sync a file:", err)
			}
			closeFile(file)
		}
	}()

	for res := range input {
		u, err := url.Parse(res.Request.Url)
		if err != nil {
			log.Fatalln("cannot parse a URL:", err)
		}

		file, ok := fileMap[u.Host]
		if !ok {
			path := filepath.Join(dir, u.Host+".csv")
			file, err = os.Create(path)
			if err != nil {
				log.Fatalln("cannot create a new file:", err)
			}
			fileMap[u.Host] = file
			writeLine(file, "url,response\n")
		}

		line := fmt.Sprintf("%v,%v\n", res.Request.Url, res.Response)
		writeLine(file, line)
	}
}

func writeLine(file *os.File, line string) {
	_, err := file.WriteString(line)
	if err != nil {
		log.Fatalln("cannot write into a file:", err)
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalln("cannot close a file:", err)
	}
}
