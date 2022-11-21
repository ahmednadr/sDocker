package operations

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
func ExtractImage(path string, ContainerID string) {
	_, checkErr := os.Open("./containers")
	if checkErr != nil {
		os.Mkdir("./containers", 0777)
	}
	os.Mkdir("./containers/"+ContainerID, 0777)
	extract := exec.Command("tar", "-xvf", path, "-C", "./containers/"+ContainerID)
	err := extract.Run()
	if err != nil {
		panic(err)
	}
}

func extractTarForBuild(baseImage string, name string) {
	_, checkErr := os.Open("./images/tmp")
	if checkErr != nil {
		os.Mkdir("./images/tmp", 0777)
	}
	os.Mkdir("./images/tmp/"+name, 0777)
	extract := exec.Command("tar", "-xvf", "./images/"+baseImage+".tar.gz", "-C", "./images/tmp/"+name)
	err := extract.Run()
	if err != nil {
		panic(err)
	}
}

func BuildNewNs() {
	buildRunCmd := exec.Command("/proc/self/exe", append([]string{"buildinternal"}, os.Args[2:]...)...)
	buildRunCmd.Stdout = os.Stdout
	buildRunCmd.Stderr = os.Stderr
	buildRunCmd.Stdin = os.Stdin

	buildRunCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	buildRunCmd.Run()
}

func inContainerThread(c chan []string) {
	fmt.Printf("running %v %d\n", os.Args[1:], os.Getpid())

	syscall.Chroot("./images/tmp/" + os.Args[2])
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("/proc", 0)

	for bashCmd := range c {

		fmt.Printf("running %v %d\n", bashCmd, os.Getpid())

		cmd := exec.Command(bashCmd[0], bashCmd[1:]...)
		// cmd := exec.Command("echo", os.Getenv("PATH"))

		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		err := cmd.Run()

		if err != nil {
			panic(err)
		}
	}
}

func Build(buildFilePath string, newImageName string) {
	f, err := os.Open(buildFilePath + "/sDockerfile")
	if err != nil {
		panic(err)
	}

	lines := bufio.NewScanner(f)
	var parsedFile [][]string
	for lines.Scan() {
		line := lines.Text()
		trimedline := strings.TrimSpace(line)
		// TODO: reduce time complxicity by picking the first word
		// rather than splitting
		args := strings.Split(trimedline, " ")
		parsedFile = append(parsedFile, args)
	}

	c := make(chan []string)

	for _, cmd := range parsedFile {
		switch cmd[0] {
		case "FROM":
			{
				extractTarForBuild(cmd[1], newImageName)
				go inContainerThread(c)
			}
		case "RUN":
			c <- cmd[1:]
		}
	}
	close(c)
}
