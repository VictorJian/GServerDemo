package Setting

import "GSFH/Models"

func SettingAdmin()  {

	if Models.FindAdminAutoID() > 0{
		Models.CreateAdminAutoId()
	}

}