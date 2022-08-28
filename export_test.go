package shuffler


// SetAfterSave replaces afterSave module-level fun with custom callback.
// For use in tests.
func SetAfterSave(f func(*Project)error){
	afterSave=f
}
