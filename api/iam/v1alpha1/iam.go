package v1alpha1

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
