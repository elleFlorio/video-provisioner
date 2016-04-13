package load

import (
	"bytes"
	"encoding/csv"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	res "github.com/elleFlorio/video-provisioner/resources"
)

type load map[float64]float64

const c_PROFILES_PATH = "profiles/"
const c_FILE_EXT = ".csv"

var (
	profiles      map[string]load
	probabilities map[string]float64
	rnd           *rand.Rand
)

func init() {
	profiles = make(map[string]load)
	probabilities = make(map[string]float64)
	source := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(source)
}

func ReadProbabilities(values []string) {
	log.Println("Reading profiles probabilities...")
	defer log.Println("Done")

	var err error
	var profProb []string
	var prob float64
	probCheck := 0.0
	for _, value := range values {
		profProb = strings.Split(value, ":")
		if len(profProb) == 2 {
			prob, err = strconv.ParseFloat(profProb[1], 64)
			if err != nil {
				log.Println("Error reading probability. Set to 0.0. profile: " + profProb[0])
				prob = 0.0
			}
		} else {
			prob = 1.0
		}
		probCheck += prob

		probabilities[profProb[0]] = prob
		log.Printf("Read profile %s having probability %f", profProb[0], prob)
	}

	if probCheck != 1.0 {
		log.Fatalln("Error: the sum of profiles probabilities should be 1.0")
	}
}

func ReadProfiles(files []string) {
	log.Println("Reading load profiles from resources...")
	defer log.Println("Done")

	for _, fileName := range files {
		log.Println("Reading load profile " + fileName)
		var err error

		file, err := res.Asset(c_PROFILES_PATH + fileName + c_FILE_EXT)
		if err != nil {
			log.Fatalln("Error getting profile " + fileName)
		}

		reader := csv.NewReader(bytes.NewReader(file))
		profile_raw, err := reader.ReadAll()
		if err != nil {
			log.Fatalln("Error reading load profile " + fileName)
		}

		profile_load := make(load)
		var prob float64
		var time float64
		for line, fields := range profile_raw {
			prob, err = strconv.ParseFloat(fields[1], 64)
			time, err = strconv.ParseFloat(fields[0], 64)
			if err != nil {
				log.Println("Error reading line " + string(line))
			}
			profile_load[time] = prob
		}

		profiles[fileName] = profile_load
	}
}

func GetProfilesNames() []string {
	names := []string{}
	for name, _ := range probabilities {
		names = append(names, name)
	}

	return names
}

func GetLoad() float64 {
	cur_load := getLoadFromProfiles(probabilities)
	return getValueFromLoad(cur_load)
}

func getLoadFromProfiles(values map[string]float64) load {
	if len(values) == 1 {
		for name, _ := range values {
			return profiles[name]
		}
	}

	profile := extractProfileFromProbabilities(values)
	return profiles[profile]

}

func getValueFromLoad(values load) float64 {
	if len(values) == 1 {
		for _, value := range values {
			return value
		}
	}

	return extractValueFromLoad(values)
}

func extractProfileFromProbabilities(probabilities map[string]float64) string {
	p := rnd.Float64()
	probSum := 0.0
	for value, prob := range probabilities {
		probSum += prob
		if p <= probSum {
			return value
		}
	}

	log.Println("Unable to extract load from profile probabilities")

	return ""
}

func extractValueFromLoad(probabilities load) float64 {
	p := rnd.Float64()
	for value, prob := range probabilities {
		if p <= prob {
			return value
		}
	}

	log.Println("Unable to extract value from load probabilities")

	return 0.0
}
