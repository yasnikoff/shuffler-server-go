package shuffler

func init() {
	afterSave = func(pr *Project) error {
		/*r := pr.Repo()
		if r != nil {
			return r.UpdateIndexForPath(pr.Path())
		}*/
		return nil
	}
}

