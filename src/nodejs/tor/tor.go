package tor
import(
	"github.com/digantaTheProgrammer/NodeTor/src/apt"
	"github.com/cloudfoundry/libbuildpack"
)

func InstallTor(s *apt.Supplier,Log *libbuildpack.Logger) error {
	Log.Info("Installing Tor.....")
	installDir := filepath.Join(s.Stager.DepDir(),"tor")
	err,options,doptions,archiveDir :=apt.AptSetup(installDir)
	if(err!=nil){
	 	return err
	 }

	err := apt.AptUpdate(options,s.Command)
	if(err!=nil){
	 	return err
	 }

	err := apt.DownloadRepoPackage(doptions,"tor",s.Command)
	if(err!=nil){
	 	return err
	 }

	err := apt.InstallPackages(archiveDir,installDir,s.Command)
	if(err!=nil){
	 	return err
	 }

	err := apt.LinkPackages(installDir,s.Stager)
	if(err!=nil){
	 	return err
	 }

	Log.Info("Tor installed!!!")
	torscript:=`export TOR_PORT_1=58`
	s.Stager.WriteProfileD("tor.sh",torscript)
	return nil
}