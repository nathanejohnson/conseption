nodefflag package

This package gives you a "no default" variant of all of the flag package types.  This means that if the flag is not passed, the value returned from the NoDef* calls or passed into NoDef*Var calls will not be set.  In order to allow this to work, we have to use double pointers everywhere.  A bit clunky but this allows us, for instance, have a config file that can set parameters, then have command line parameters override those.


