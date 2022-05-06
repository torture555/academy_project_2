package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type descAscOption int // initializing type and const descending, ascending

const (
	descending descAscOption = 1
	ascending                = 2
)

type fileInfo struct {
	path    string
	size    int64
	hashSum string
} // initializing struct fileInfo with attributes a path to file, size a file, hashSum a file

func (m *fileInfo) calculateHashSum() error {
	file, err := os.Open(m.path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println(err)
		return err
	}
	tmpHashSum := hash.Sum(nil)
	strHashSum := hex.EncodeToString(tmpHashSum)
	m.hashSum = strHashSum
	return nil
} // func to calculate hashSum a file by path

func newSortOption(val int) descAscOption {
	if val == 1 {
		return descending
	} else if val == 2 {
		return ascending
	}
	return 0
} // constructor type descAscOption

func newfileInfo(path string, size int64) fileInfo {
	return fileInfo{path: path, size: size, hashSum: ""}
} // constructor type fileInfo without calculate hashSum

func getPathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		err := errors.New("directory is not specified")
		return "", err
	} else {
		return os.Args[1], nil
	}
} // get input path from arguments

func getInputFormatFile() string {
	fmt.Println("Enter file format:")
	return getScanValue()
} // get format files from input value

func getSortOption() (descAscOption, error) {
	fmt.Println("Size sorting options:\n1. Descending\n2. Ascending\n\nEnter a sorting option:")
	inputStr := ""
	for inputStr = getScanValue(); inputStr != "1" && inputStr != "2"; inputStr = getScanValue() {
		fmt.Println("Wrong option")
	}
	res, err := strconv.Atoi(inputStr)
	if err != nil {
		return 0, err
	}
	return newSortOption(res), nil
} // get sorting option from input value

func getCheckForDuplicates() bool {
	fmt.Println("\nCheck for duplicates?")
	for {
		res := getScanValue()
		if res == "yes" {
			return true
		} else if res == "no" {
			return false
		}
		fmt.Println("Wrong option")
	}
} // get check value for output duplicates by hashSum from input value

func getScanValue() string {
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text()
} // scan and return user input value

func walkFilesByPath(path string, sliceFilesInfo *[]fileInfo, userFormatFile string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, userFormatFile) {
			*sliceFilesInfo = append(*sliceFilesInfo, newfileInfo(path, info.Size()))
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
} // walk files by path and append in slice with objects struct fileInfo

func cutSliceHashSumInfoFiles(sliceFilesInfo *[]fileInfo) []string {
	var slice []string
	for _, obj := range *sliceFilesInfo {
		slice = append(slice, obj.hashSum)
	}
	return slice
} // cut slice with hashSums from slice objects fileInfo

func cutSliceSizeInfoFiles(sliceFilesInfo *[]fileInfo) []int64 {
	var slice []int64
	for _, obj := range *sliceFilesInfo {
		slice = append(slice, obj.size)
	}
	return slice
} // cut slice with sizes files from slice objects fileInfo

func sortingSliceFileInfoBySize(slice *[]fileInfo, sortOption descAscOption) {
	tmpSlice := *slice
	if sortOption == descending {
		sort.Slice(tmpSlice, func(i, j int) bool { return tmpSlice[i].size > tmpSlice[j].size })
	} else if sortOption == ascending {
		sort.Slice(tmpSlice, func(i, j int) bool { return tmpSlice[i].size < tmpSlice[j].size })
	}
} // sorting slice with objects fileInfo by size

func groupSlice[t int64 | string](slice *[]t) {
	tmpSlice := *slice
	for i := 0; i < len(tmpSlice)-1; i++ {
		if tmpSlice[i] == tmpSlice[i+1] {
			tmpSlice = append(tmpSlice[:i], tmpSlice[i+1:]...)
			//			end = false
			i--
		}
	}
	*slice = tmpSlice
} // grouping slice with sizes/hashSums values

func groupDuplicationsHashSumsFromSlice(slice *[]string) {
	tmpSlice := *slice
	var resSlice []string
	for i := 0; i < len(tmpSlice)-1; i++ {
		if tmpSlice[i] == tmpSlice[i+1] {
			tmpSlice = append(tmpSlice[:i], tmpSlice[i+1:]...)
			resSlice = append(resSlice, tmpSlice[i])
			i--
		}
	}
	groupSlice(&resSlice)
	*slice = resSlice
} // get slice with duplications hashSums and grouping

func getGroupSliceMapSizeFiles(sliceFilesInfo *[]fileInfo) []int64 {
	res := cutSliceSizeInfoFiles(sliceFilesInfo)
	groupSlice(&res)
	return res
} // get grouping and sorting slice with sizes files from slice objects fileInfo

func getGroupDuplicationSliceHashFiles(sliceFilesInfo *[]fileInfo) []string {
	slice := cutSliceHashSumInfoFiles(sliceFilesInfo)
	groupDuplicationsHashSumsFromSlice(&slice)
	return slice
} // get grouping and sorting slice with duplications hashSums from slice objets fileInfo

func calculateSliceFilesInfo(sliceFilesInfo *[]fileInfo) error {
	tmpSlice := *sliceFilesInfo
	for i := 0; i < len(tmpSlice); i++ {
		err := tmpSlice[i].calculateHashSum()
		if err != nil {
			return err
		}
	}
	*sliceFilesInfo = tmpSlice
	return nil
} // calculate hashSum for each file by path from slice objects fileInfo

func printOutputSizeSorting(sliceFilesInfo *[]fileInfo, groupSize *[]int64) {
	for _, intSize := range *groupSize {
		fmt.Println("\n", intSize, "bytes")
		for _, obj := range *sliceFilesInfo {
			if obj.size == intSize {
				fmt.Println(obj.path)
			}
		}
	}
} // print file paths by grouped size

func printOutputHashDuplication(sliceFilesInfo *[]fileInfo, groupHash *[]string, sortOption descAscOption) []string {
	var menuForDelete []string
	var tmpSLiceFilesInfo []fileInfo

	for _, obj := range *sliceFilesInfo {
		for _, strHashSum := range *groupHash {
			if strHashSum == obj.hashSum {
				tmpSLiceFilesInfo = append(tmpSLiceFilesInfo, obj)
			}
		}
	}

	sortingSliceFileInfoBySize(&tmpSLiceFilesInfo, sortOption)

	str := "\n"
	count := 1
	for i := 0; i < len(tmpSLiceFilesInfo); i++ {
		menuForDelete = append(menuForDelete, tmpSLiceFilesInfo[i].path)
		str = str + strconv.Itoa(count) + ". " + tmpSLiceFilesInfo[i].path + "\n"
		if i == len(tmpSLiceFilesInfo)-1 {
			str = "\nHash: " + tmpSLiceFilesInfo[i].hashSum + str
			str = strconv.FormatInt(tmpSLiceFilesInfo[i].size, 10) + " bytes" + str
			fmt.Println(str)
			break
		}

		if tmpSLiceFilesInfo[i].hashSum != tmpSLiceFilesInfo[i+1].hashSum {
			str = "\nHash: " + tmpSLiceFilesInfo[i].hashSum + str
		}
		if tmpSLiceFilesInfo[i].size != tmpSLiceFilesInfo[i+1].size {
			str = strconv.FormatInt(tmpSLiceFilesInfo[i].size, 10) + " bytes" + str
			fmt.Println(str)
			str = "\n"
		}
		count++
	}

	return menuForDelete
} // print path files with duplications hashSums by grouped size

func main() {

	path, err := getPathFromArgs() // initializing root path
	if err != nil {
		fmt.Println("Directory is not specified")
		return
	}

	userFormatFile := getInputFormatFile() // get format files for selection

	var sliceFilesInfo []fileInfo // declare slice with objects fileInfo

	err = walkFilesByPath(path, &sliceFilesInfo, userFormatFile) // walk to file by root path and append objects in slice with fileInfo
	if err != nil {
		return
	}

	sortOption, err := getSortOption() // get sort option
	if err != nil {
		return
	}

	sortingSliceFileInfoBySize(&sliceFilesInfo, sortOption) // sorting slice with fileInfo by sort option

	groupsSize := getGroupSliceMapSizeFiles(&sliceFilesInfo) // get slice with grouped and sorted sizes from slice with objects fileInfo

	printOutputSizeSorting(&sliceFilesInfo, &groupsSize) // print file paths by grouped size

	if !getCheckForDuplicates() { // if not output duplications files by hashSums
		return
	}

	err = calculateSliceFilesInfo(&sliceFilesInfo) // calculate each files by slice with objects fileInfo
	if err != nil {
		return
	}

	groupsHash := getGroupDuplicationSliceHashFiles(&sliceFilesInfo) // get slice with grouped and sorted hashSums from slice with objects fileInfo

	printOutputHashDuplication(&sliceFilesInfo, &groupsHash, sortOption) // print path files with duplications hashSums by grouped size

}
