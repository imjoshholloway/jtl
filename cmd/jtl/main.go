package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/imjoshholloway/jtl"
)

var (
	specFile   = flag.String("spec", "", "--spec=./path/to/spec.yaml")
	sourceFile = flag.String("source", "", "--source=./path/to/source.json")
)

func specFromReader(r io.Reader) (jtl.Spec, error) {
	decoder := yaml.NewDecoder(r)

	spec := jtl.Spec{}

	for {
		s := jtl.Spec{}
		if err := decoder.Decode(&s); err != nil {
			// keep decoding specs until we hit EOF
			if err == io.EOF {
				break
			}

			return spec, err
		}

		spec.Specs = append(spec.Specs, s)
	}

	return spec, nil
}

func inputFromReader(r io.Reader) (interface{}, error) {

	var out interface{}

	dec := json.NewDecoder(r)
	for {
		err := dec.Decode(&out)
		if err == io.EOF {
			break
		}

		return out, err
	}

	return out, nil
}

func main() {

	flag.Parse()

	rawSpec, err := os.Open(*specFile)
	if err != nil {
		log.Fatal("Unable to open spec file", err)
	}

	spec, err := specFromReader(rawSpec)

	var inputSource io.Reader

	// default to read the data from stdin
	inputSource = os.Stdin

	if *sourceFile != "" {
		var err error
		inputSource, err = os.Open(*sourceFile)
		if err != nil {
			log.Fatal("Unable to open source file", err)
		}
	}

	input, err := inputFromReader(inputSource)
	if err != nil {
		log.Fatal("Unable to load source", err)
	}

	transformed := spec.Process(input)

	raw, _ := json.Marshal(transformed)

	fmt.Printf("%s", raw)
}
