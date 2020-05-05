package apt
import(
	"github.com/cloudfoundry/libbuildpack"
	"path/filepath"
)

func SingleInstall(s *Supplier,pkg string) error {
	installDir := InstallDir(s.Stager,pkg)
	err,options,doptions,archiveDir :=apt.AptSetup(s,installDir)
	if(err!=nil){
	 	return err
	 }

	err = apt.AptUpdate(options,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = apt.DownloadRepoPackage(doptions,pkg,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = apt.InstallPackages(archiveDir,installDir,s.Command)
	if(err!=nil){
	 	return err
	 }

	err = apt.LinkPackages(installDir,s.Stager)
	if(err!=nil){
	 	return err
	 }
	return nil
}