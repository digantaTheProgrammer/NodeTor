package apt
func DownloadRepoPackage(options String[],pkg String,command Command) error
{
	aptArgs :=append(options,pkg)
	out, err := command.Output("/", "apt-get", aptArgs...)
	if err != nil {
		return fmt.Errorf("failed apt-get install %s\n\n%s", out, err)
	}
	return nil;
}
