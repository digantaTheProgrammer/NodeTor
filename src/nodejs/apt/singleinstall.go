package apt

func SingleInstall(s *Supplier,pkg string) error {
	installDir := InstallDir(s.Stager,pkg)
	err,options,doptions,archiveDir :=AptSetup(s,installDir)
	if(err!=nil){
	 	return err
	 }

	err = AptUpdate(options,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = DownloadRepoPackage(doptions,pkg,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = InstallPackages(archiveDir,installDir,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = LinkPackages(installDir,s.Stager)
	if(err!=nil){
	 	return err
	 }
	return nil
}