ifit:  A simple tool for manipulating source code with if and substitution
==========================================================================

 [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/pschlump/Go-FTL/master/LICENSE)


*ifit* pronounced *if* *it* - not pronounced *I* *FIT* (not like *I* had a *fit* one day and went off and wrote this -- and not like 
*I* am *FIT* and in good shape).  

*ifit* is a little like GNU m4 or the C preprocessor.  It allows you to have source level `if` statements and substitute
for values to create a configured program.  The indented source is HTML, CSS and JavaScript.  If it were to be used for
other languages then additional syntax might be needed.

This program was prompted by having a project with 98% the same set of source code that really couldn't have the last 2%
removed at run time.  iOS/Safari required a file that made Android/Chrome break.  Thus a tool was born. 

Command Line Aruments
---------------------

Argument | Long       | Description
:---:    | ---------- | --------------------------------------------------------------------
`-i`     | `--input`  | Input file
`-o`     | `--output` | Output file
`-s`     | `--sub`    | Substitution file in JSON
`-m`     | `--mode`   | The subset of the configuration file that will be used, "prod", "dev" etc.
`-D`     | `--debug`  | Debug flag true/false

Anny additional arguments are takes as turning on that defined section of code.  For example:

```
	<!-- !! if IOS_SECTION !! -->
		this section of HTML is for my IOS version only!
	<!-- !! else IOS_SECTION !! -->
		this section of HTML is for my non-IOS version (NOT IOS)!
	<!-- !! end IOS_SECTION !! -->
```

or

```
	<!-- !! if IOS_SECTION !! -->
		this section of HTML is for my IOS version only!
	<!-- !! end IOS_SECTION !! -->
```

and you run

```
	$ ifit -i input.html -o output.html IOS_SECTION OtherSection
```

then you will include the IOS_SECTION.  If you run

```
	$ ifit -i input.html -o output.html ANDROID_SECTION OtherSection
```

you will not include it.

Substitution Values
-------------------

The -s option allows you to read in a JSON file of substitution values.  These are also taken to be
sections that you would want to have turned on.

Example:

```
{
	"__comment__": {
			"line1":"Configuration for iOS in the ./iOS directory"
			,"line2":"Copyright (C) Philip Schlump, 2016."
			,"line3":"MIT Licnesed."
		}
	,"dev":
		{
			"iOS": "on"
			,"Mobile": "on"
			,"DevMode": "on"
			,"requestMode":"jsonp"
			,"pathToAesSrpServer":"http://localhost:3118/api/"
		}
	,"test":
		{
			"iOS": "on"
			,"Mobile": "on"
			,"DevMode": "on"
			,"requestMode":"jsonp"
			,"pathToAesSrpServer":"http://localhost:3118/api/"
		}
	,"prod":
		{
			"iOS": "on"
			,"Mobile": "on"
			,"ProdMode": "on"
			,"requestMode":"POST"
			,"pathToAesSrpServer":"https://www.go-ftp.com/api/"
		}
}
```

will substitute `$$iOS$$` for `on`.  It will also turn on the section `<!-- !! if iOS !! -->`.

The flags on the command line override the values in the `-s`/`--sub` file.

Syntax
------

HTML

```
<!-- !! if NameA !! -->
<div id="A"></div>
<!-- !! end NameA !! -->
```

or

```
<!-- !! if NameB !! 
<div id="B"></div>
!! end NameB !! -->
```

JavaScript

```
// !! if NameA !!
var A = 12;
// !! end NameA !! 
```

CSS

```
/* !! if NameC !! */
.classC {
}
/* !! end NameC !!  */
```

or

```
/* !! if NameD !!
.classD {
}
!! end NameD !!  */
```

Predefined Values
-----------------

The following predefined values are set.

Name                   | Description
-----------------------|-------------------------------------------------------------------------------------------------
`$$__FILE__$$`         | Current Input File Name
`$$__LINE__$$`         | Current Line Number
`$$__DATE__$$`         | Date in ISO, YYYY/MM/DD format
`$$__TIME__$$`         | 24 hour time in ISO, HH24:MI:SS format
`$$__TSTAMP__$$`       | Timestamp in RFC3339 format
`$$__Mode__$$`         | Current mode from the `-m` flag
`$$__Output__$$`       | Current output file name
`$$__TRUE_ITEMS__$$`   | Items that are defined to be true via command line or via input file.  In sorted order.
`$$__OPENED_FILES__$$` | Name of all the files that have been opened after adding the search path.
`$$__PATH__$$`         | Current search path.

Predefined Operators
--------------------

The following predefined operators:

Name           | Description
---------------|----------------------------------------------------------------------------------------------------------------- 
`define`       | Define a name to be a value.  The same as placing on the command line NAME=Value or in the `-m` file.
`undef`        | Removes a defined value
`set_path`     | Sets the search path for include files
`include`      | Include a file in the middle of the output.
`include_once` | Include a file, if it has not already been included, in the middle of the output.
`if`           | If the item is defined then include the following section up to a matching `end` or `else`.
`end`          | end of an if
`else`         | you know.  In the middle of an if with the reversed logic.

Please Note
-----------

Tests are in a `Makefile` and run by 

```
	$ make test0
	$ make test1
	$ make test2
	...
	$ make test11
```

You should see *PASS* at the end of each successful test.

LICENSE
-------

MIT Licensed -  See LICENSE file.

Author
------

By Philip Schlump.

TODO
----

### Documentation

1. Add examples in this file
1. Have somebody edit the documentation (this file) to improve it

### Code Cleanup

1. Add tests for ./stk
1. Add tests for ./fstk

### Add Features

1. Add elseif - Simple
1. Add expressions ( iOS || Android ), ( !XDesktop ) - Medium
1. Add `$$__noeol__$$` - Hm....
1. Allow set of Start/End marker for tokens instead of fixed `$$` and `$$` - Medium
1. Add in parameters for defined stuff and quotes and non-eval expressions like m4 - Hard

