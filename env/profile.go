package env

type Profile int8

const (
	_invalid = "invalid"
	_local   = "local"
	_sandbox = "sandbox"
	_dev     = "dev"
	_staging = "staging"
	_qa      = "qa"
	_prod    = "prod"
)

const (
	Invalid Profile = iota - 1
	Local
	Sandbox
	Dev
	Staging
	QA
	Prod
)

func (p Profile) String() string {
	switch p {
	case Local:
		return _local
	case Dev:
		return _dev
	case Sandbox:
		return _sandbox
	case Staging:
		return _staging
	case QA:
		return _qa
	case Prod:
		return _prod
	case Invalid:
		return _invalid
	default:
		return _invalid
	}
}

// IsValid check the profile is valid
func (p Profile) IsValid() bool {
	return Local <= p && p <= Prod
}

func stringToProfile(str string) Profile {
	switch str {
	case _local:
		return Local
	case _sandbox:
		return Sandbox
	case _dev:
		return Dev
	case _staging:
		return Staging
	case _qa:
		return QA
	case _prod:
		return Prod
	default:
		return Invalid
	}
}
