package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"github.com/FiloSottile/powersoftau/powersoftau"
)

func main() {
	challengeFile := flag.String("challenge", "./challenge", "path to the challenge file")
	responseFile := flag.String("response", "./response", "path to the response file")
	pprof := flag.Bool("pprof", false, "run a profiling server; use ONLY FOR DEBUGGING")
	flag.Parse()

	if *pprof {
		go http.ListenAndServe("localhost:6060", nil)
	}

	log.Printf("Reading challenge...\n")
	ch, err := powersoftau.ReadChallenge(*challengeFile)
	if err != nil {
		log.Fatalf("Failed to read the challenge: %v\n", err)
	}

	log.Printf("Starting computation...\n")
	ch.Compute(runtime.NumCPU())

	log.Printf("Writing response...\n")
	hash, err := powersoftau.WriteResponse(*responseFile, ch)
	if err != nil {
		log.Fatalf("Failed to write the response: %v\n", err)
	}

	log.Printf("Done!\n\nYour contribution has been written to `%s`\n\nThe BLAKE2b hash of `%s` is:\n", *responseFile, *responseFile)
	for i := 0; i < 4; i++ {
		fmt.Printf("\t")
		for k := 0; k < 4; k++ {
			fmt.Printf("%x ", hash[i*4*4+k*4:i*4*4+k*4+4])
		}
		fmt.Printf("\n")
	}
}
