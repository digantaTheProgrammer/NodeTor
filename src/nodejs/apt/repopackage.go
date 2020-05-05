package apt
import(
	"fmt"
)
func DownloadRepoPackage(options []string,pkg string,command Command) error{
	aptArgs :=append(options,pkg)
	out, err := command.Output("/", "apt-get", aptArgs...)
	if err != nil {
		return fmt.Errorf("failed apt-get install %s\n\n%s", out, err)
	}
	return nil;
}
