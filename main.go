package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type pair struct{
	hash, path string
}

type fileList []string
type result map[string]fileList

func getHash(path string)pair{
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		log.Println(err)
	}

	return pair{fmt.Sprintf("%x", hash.Sum(nil)), path}
}

func collectHashes(p <-chan pair, r chan<- result){
	hashes := make(result)

	for val := range p{
		hashes[val.hash] = append(hashes[val.hash], val.path)
	}
	r <- hashes
}

func processFiles(path <-chan string, pair chan<-pair, done chan<-bool)  {
	for p := range path{
		pair <- getHash(p)
	}
	done <- true
}

func searchTree(dir string, p chan<-string)error{

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		fileInfo, err := d.Info();
		if err != nil {
			log.Println(err)
		}

		if d.Type().IsRegular() && fileInfo.Size() > 0 {
			p <- path
		}
		return nil
	})
	return err
}

func run(dir string)result{

	workers := 2*runtime.GOMAXPROCS(0)
	path := make(chan string)
	pair := make(chan pair)
	result := make(chan result)
	done := make(chan bool)

	for i := 0; i < workers; i++ {
		go processFiles(path, pair, done)
	}

	go collectHashes(pair, result)

	if err := searchTree(dir, path); err != nil {
		log.Println(err)
	}

	close(path);

	for i := 0; i < workers; i++ {
		<-done
	}

	close(pair)

	return <-result
}

func main(){
	totalFiles := 0;
	start := time.Now()

	pathPtr := flag.String("path", "./", "Path to search in Directory")
	flag.Parse()
	fmt.Println(*pathPtr)

	if hashes := run(*pathPtr); hashes != nil{
		l := time.Since(start).Round(time.Millisecond)
		for hash, files := range hashes{
			if(len(files) >= 1){
				totalFiles += len(files)
				fmt.Println("Hash is....", hash)
				for i, file := range files{
					fmt.Println(i, " ", file)
				}
			}
		}
		fmt.Println("Total files: ", totalFiles)
		fmt.Println("Time taken ", l)
		timeTakenPerFile := float64(l.Milliseconds())/float64(totalFiles)
		fmt.Println(timeTakenPerFile, " files/mili")
	}
}