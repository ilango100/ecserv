package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func init() {

	if ok, err := loadSettings(); ok {
		//Settings file done

	} else if os.IsNotExist(err) {
		//First start, run wizard
		firstStart()
	} else {
		log.Fatalln("Error reading settings file...", err)
	}

	//finally settings will be loaded
	if invalidSettings() {
		log.Fatalf("Invalid Settings...\nCorrect the settings file %s and try again...", isfile)
	}
}

func firstStart() {

	if s {
		log.Fatalln("Silent mode without settings file invalid")
	}

	var c int

	fmt.Print("***Welcome to EcServ webserver!***\n\nWe need to set settings now...\n\n1. Express settings (default)\n2. Custom (experts)\n\nYour choice...")

	if _, err := fmt.Scan(&c); err != nil || c > 2 || c < 1 {
		fmt.Println("Wrong option, defaulting to Express setting...")
		c = 1
	}

	//Express settings
	if c == 1 {

		//root dir
		if homedir := os.Getenv("HOME"); homedir != "" {
			set.Root = homedir + "/EcServ"
		} else if os.Getenv("OS") == "Windows_NT" {
			set.Root = os.Getenv("USERPROFILE") + "\\EcServ"
		}
		fmt.Println("\nDefault root directory...", set.Root)

		//email
		getMail()

		//Certificate dir
		fmt.Println("Default certificate directory...cert")
		set.Cert = "cert"

		//Domain
		getDomain()

		//Error log file
		set.ErrLog = "errors.log"

		//settings file
		isfile = "ecset"

		//Custom settings
	} else {
		//root dir
		fmt.Print("Enter root directory...")
		fmt.Scan(&set.Root)

		//email
		getMail()

		//Certificate dir
		fmt.Print("Enter certificate directory (d - Default cert)...")
		fmt.Scan(&set.Cert)
		if set.Cert == "d" {
			set.Cert = "cert"
		}

		//Domain
		getDomain()

		//Error log file
		fmt.Print("Enter the error log file name (d - Default errors.log)...")
		fmt.Scan(&set.ErrLog)
		if set.ErrLog == "d" {
			set.ErrLog = "errors.log"
		}

		//settings file
		if isfile == "ecset" {
			fmt.Print("\nEnter filename to write settings (d - Default ecset)...")
			fmt.Scan(&isfile)
			if isfile == "d" {
				isfile = "ecset"
			}
		}

	}

	fmt.Printf("Settings over... Writing settings to file %s...\n", isfile)

	if iset, err := os.Create(isfile); err == nil {
		ency := json.NewEncoder(iset)
		ency.SetIndent("", " ")
		if err = ency.Encode(set); err != nil {
			log.Fatalln("Error writing settings file...", err)
		}
		iset.Close()
	} else {
		log.Fatalln("Error writing settings file...", err)
	}
}

func wrongDir(str string) bool {
	if inf, err := os.Stat(str); err != nil {
		if os.IsNotExist(err) {
			//not silent, so ask first
			if !s {
				fmt.Printf("Directory %s non-existing...\nCreate now? (Enter 1) ", str)
				var c int
				fmt.Scan(&c)
				if c != 1 {
					log.Fatal("Try again...", c)
				}
			}
			//Try to creat dir
			if err = os.Mkdir(str, 0755); err != nil {
				log.Fatalln("Error creating directory", str, err)
			}

		} else {
			log.Fatalln("Error reading the directory...", err)
		}
	} else if inf.IsDir() {
		//Directory valid!

	} else {
		//Exists but not a directory
		log.Fatalln("Root exists but not a directory...")
	}
	return false
}

func getMail() {
	fmt.Print("Enter email...")
	fmt.Scan(&set.Email)
}

func getDomain() {
	fmt.Print("Enter domain...")
	fmt.Scan(&set.Domain)
}

func loadSettings() (bool, error) {
	//try to open settings file
	if ecset, err := os.Open(isfile); err == nil {

		//if opens decode it
		setr := json.NewDecoder(ecset)
		if err = setr.Decode(&set); err != nil && err != io.EOF {
			return false, err
		} else {
			ecset.Close()
			return true, nil
		}

	} else {
		return false, err
	}
}

func invalidSettings() (wrong bool) {
	if set.Root == "" {
		fmt.Println("Root dir not set...")
		wrong = true
	} else if wrongDir(set.Root) {
		wrong = true
	}
	if set.Email == "" {
		fmt.Println("email not set...")
		wrong = true
	} else if ok, _ := regexp.MatchString("(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)", set.Email); !ok {
		fmt.Println("email not valid")
		wrong = true
	}
	if set.Cert == "" {
		fmt.Println("Certificate folder not specified...")
		wrong = true
	} else if wrongDir(set.Cert) {
		wrong = true
	}
	if set.Domain == "" {
		fmt.Println("Domain not set...")
		wrong = true
	} else if ok, _ := regexp.MatchString("^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\\.[a-zA-Z]{2,3})$", set.Domain); !ok {
		fmt.Println("Domain name invalid...")
		wrong = true
	}
	return wrong
}
