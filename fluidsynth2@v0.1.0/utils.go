package fluidsynth2

type MenuObject struct{
	Name string
}

func (menuObject MenuObject) MenuObjectName() string {
    return menuObject.Name
}