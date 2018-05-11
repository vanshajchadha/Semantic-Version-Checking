package main

import (
	"bufio"
	"os"
	"context"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	
	m:= make(map[string]semver.Version)

	for i:=0;i<len(releases);i++ {
		if releases[i].LessThan(*minVersion)==false {
			s:= fmt.Sprintf("%d.%d",releases[i].Major,releases[i].Minor) // Storing all combinations of major and minor versions in a given format as keys
			_,ok:= m[s]
			if ok == true {
				elem:= m[s] 
				if releases[i].Patch > elem.Patch {
					m[s] = *releases[i]
				} else if releases[i].Patch == elem.Patch {
					 /*if semver.preReleaseCompare(releases[i],elem) == 1 {	// No need to implement this as we are asked to fetch simply according to the highest patch version in the question
					 	m[s] = *releases[i]
					 }*/						
				}

			} else {
				m[s] = *releases[i]
			}
		}
	}

	var temp []semver.Version

	for k := range m {
		temp = append(temp,m[k])	// Creating a temporary array which will have all the values corresponding to distinct keys
	}

	for k := range temp {
		versionSlice = append(versionSlice,&temp[k])	
	}

	for i:=0;i<len(versionSlice);i++ {
 		max := versionSlice[i]
 		pos := i

 		for j:=i+1;j<len(versionSlice);j++ {
 			if versionSlice[j].Major > max.Major {	// Using the concept of Selection Sort. We could have also used any other comparing sort mechanism
	 			max=versionSlice[j]
	 			pos=j
	 		} else if versionSlice[j].Major == max.Major && versionSlice[j].Minor > max.Minor {
	 			max=versionSlice[j]
	 			pos=j
	 		}
 		}							// We just need to compare the major and minors because each combination of major and minor versions would be unique now

 		temp:=versionSlice[i]
 		versionSlice[i]=versionSlice[pos]
 		versionSlice[pos]=temp
	}
				
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	fileAddress := os.Args[1]
	file, err1 := os.Open(fileAddress)
	var repo1 []string
	var repo2 []string
	var minver []string
	
	if(err1 == nil){
		scanner := bufio.NewScanner(file)
		first:=true							// Used to ignore the first line while reading the file as the first line contains general format of the whole file
		for scanner.Scan() {
			line := scanner.Text()
			if first == false {
				c:=0
				for i:=0;i<len(line);i++ {
					if(line[i]=='/'){
						repo1=append(repo1,line[:i])
						c=i+1
					} else if(line[i]==','){
						repo2=append(repo2,line[c:i])
						minver=append(minver,line[i+1:])
						c=0
						break
					}
				}
			}
			first=false
		}
	}

	defer file.Close()

	for i:= range repo1 {
		
		releases, _, err := client.Repositories.ListReleases(ctx, repo1[i], repo2[i], opt)
		if err != nil {
			fmt.Println(err)	// Instead of panicking, it is better to print the error so that the user can debug easily and return from the function
			return
		}
		minVersion := semver.New(minver[i])
		allReleases := make([]*semver.Version, len(releases))
		for i, release := range releases {
			versionString := *release.TagName
			if versionString[0] == 'v' {
				versionString = versionString[1:]
			}
			allReleases[i] = semver.New(versionString)
		}
		versionSlice := LatestVersions(allReleases, minVersion)

		fmt.Printf("latest versions of %s/%s: %s\n", repo1[i], repo2[i], versionSlice)
	}

}