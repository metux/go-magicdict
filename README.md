# go-magicdict
Recursive dictionary w/ yaml loading and variable substitution

The magicdict package provides loading YAML and acessing (getting and setting)
arbitrary entries by sinple Key pathes, while also supporting default values
as well as variable substition.

Key pathes
----------

A key path (or just Key) is quite like a path name, but using "::" as delimiter.
Non-scalar entries may have sub entries, that can be accessed by subkeys, similar
fike or directories within another directory in a filesystem.

Magic attributes
----------------

Keys starting with "@@" prefix have very special meaning:
Instead of accessing sub entries, providing special functions,
eg. referring to the entry's parent, etc. These also can also
be used in variable references, eg. for relative references.

Variable substitution
---------------------

Scalar entries (strings) may contain references to other
pathes, which are substituted when the entry is retrieved:

    ${foo}-${bar}
    ${one::two}

If an entry only contains exactly one reference, it returns
exactly the referred entry (keeping it's type) w/o trying
to convert it to string, thus allows referring to dicts or
lists, and accessing their subkeys - similar to symbolic
links in a filesystem.


