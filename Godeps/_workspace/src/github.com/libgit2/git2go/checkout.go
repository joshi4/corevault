package git

/*
#include <git2.h>
*/
import "C"
import (
	"os"
	"runtime"
)

type CheckoutStrategy uint

const (
	CheckoutNone                      CheckoutStrategy = C.GIT_CHECKOUT_NONE                         // Dry run, no actual updates
	CheckoutSafe                                       = C.GIT_CHECKOUT_SAFE                         // Allow safe updates that cannot overwrite uncommitted data
	CheckoutSafeCreate                                 = C.GIT_CHECKOUT_SAFE_CREATE                  // Allow safe updates plus creation of missing files
	CheckoutForce                                      = C.GIT_CHECKOUT_FORCE                        // Allow all updates to force working directory to look like index
	CheckoutAllowConflicts                             = C.GIT_CHECKOUT_ALLOW_CONFLICTS              // Allow checkout to make safe updates even if conflicts are found
	CheckoutRemoveUntracked                            = C.GIT_CHECKOUT_REMOVE_UNTRACKED             // Remove untracked files not in index (that are not ignored)
	CheckoutRemoveIgnored                              = C.GIT_CHECKOUT_REMOVE_IGNORED               // Remove ignored files not in index
	CheckotUpdateOnly                                  = C.GIT_CHECKOUT_UPDATE_ONLY                  // Only update existing files, don't create new ones
	CheckoutDontUpdateIndex                            = C.GIT_CHECKOUT_DONT_UPDATE_INDEX            // Normally checkout updates index entries as it goes; this stops that
	CheckoutNoRefresh                                  = C.GIT_CHECKOUT_NO_REFRESH                   // Don't refresh index/config/etc before doing checkout
	CheckooutDisablePathspecMatch                      = C.GIT_CHECKOUT_DISABLE_PATHSPEC_MATCH       // Treat pathspec as simple list of exact match file paths
	CheckoutSkipUnmerged                               = C.GIT_CHECKOUT_SKIP_UNMERGED                // Allow checkout to skip unmerged files (NOT IMPLEMENTED)
	CheckoutUserOurs                                   = C.GIT_CHECKOUT_USE_OURS                     // For unmerged files, checkout stage 2 from index (NOT IMPLEMENTED)
	CheckoutUseTheirs                                  = C.GIT_CHECKOUT_USE_THEIRS                   // For unmerged files, checkout stage 3 from index (NOT IMPLEMENTED)
	CheckoutUpdateSubmodules                           = C.GIT_CHECKOUT_UPDATE_SUBMODULES            // Recursively checkout submodules with same options (NOT IMPLEMENTED)
	CheckoutUpdateSubmodulesIfChanged                  = C.GIT_CHECKOUT_UPDATE_SUBMODULES_IF_CHANGED // Recursively checkout submodules if HEAD moved in super repo (NOT IMPLEMENTED)
)

type CheckoutOpts struct {
	Strategy       CheckoutStrategy // Default will be a dry run
	DisableFilters bool             // Don't apply filters like CRLF conversion
	DirMode        os.FileMode      // Default is 0755
	FileMode       os.FileMode      // Default is 0644 or 0755 as dictated by blob
	FileOpenFlags  int              // Default is O_CREAT | O_TRUNC | O_WRONLY
}

func (opts *CheckoutOpts) toC() *C.git_checkout_options {
	if opts == nil {
		return nil
	}
	c := C.git_checkout_options{}
	populateCheckoutOpts(&c, opts)
	return &c
}

// Convert the CheckoutOpts struct to the corresponding
// C-struct. Returns a pointer to ptr, or nil if opts is nil, in order
// to help with what to pass.
func populateCheckoutOpts(ptr *C.git_checkout_options, opts *CheckoutOpts) *C.git_checkout_options {
	if opts == nil {
		return nil
	}

	C.git_checkout_init_options(ptr, 1)
	ptr.checkout_strategy = C.uint(opts.Strategy)
	ptr.disable_filters = cbool(opts.DisableFilters)
	ptr.dir_mode = C.uint(opts.DirMode.Perm())
	ptr.file_mode = C.uint(opts.FileMode.Perm())

	return ptr
}

// Updates files in the index and the working tree to match the content of
// the commit pointed at by HEAD. opts may be nil.
func (v *Repository) CheckoutHead(opts *CheckoutOpts) error {
	var copts C.git_checkout_options

	ptr := populateCheckoutOpts(&copts, opts)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_checkout_head(v.ptr, ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}

// Updates files in the working tree to match the content of the given
// index. If index is nil, the repository's index will be used. opts
// may be nil.
func (v *Repository) CheckoutIndex(index *Index, opts *CheckoutOpts) error {
	var copts C.git_checkout_options
	ptr := populateCheckoutOpts(&copts, opts)

	var iptr *C.git_index = nil
	if index != nil {
		iptr = index.ptr
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ret := C.git_checkout_index(v.ptr, iptr, ptr)
	if ret < 0 {
		return MakeGitError(ret)
	}

	return nil
}
