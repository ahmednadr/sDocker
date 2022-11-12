package operations

import (
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func GenerateUID(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345678"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes)-1)]
	}
	return string(b)

}

// tar - xvzf <tar file name with .tar.gz extension>
func Extract(path string, ContainerID string) {
	_, checkErr := os.Open("./containers")
	if checkErr != nil {
		os.Mkdir("./containers", 0755)
	}
	os.Mkdir("./containers/"+ContainerID, 0755)
	extract := exec.Command("tar", "-xvf", path, "-C", "./containers/"+ContainerID)
	err := extract.Run()
	if err != nil {
		panic(err)
	}
}
