ifit:  A simple tool for manipulating source code with if and substitution
==========================================================================

*ifit* is a little like GNU m4 or the C preprocessor.  It allows you to have source level `if` statements and substitute
for values to create a configured program.  The indented source is HTML, CSS and JavaScript.  If it were to be used for
other languages then additional syntax might be needed.

This program was prompted by having a 98% the same set of source code that really couldn't have the last 2%
at run time.   Thus a tool was born.

Command Line Aruments
---------------------

Argument | Long Argument| Description
:---: | --- | ---
`-i` | `--input` | Input file
`-o` | `--output` | Output file
`-s` | `--sub` | Substitution file in JSON
`-D` | `--debug` | Debug flag true/false

Anny additional arguments are takes as turning on that defined section of code.  For example:

	<!-- !! if IOS_SECTION !! -->
		this section of HTML is for my IOS version only!
	<!-- !! end IOS_SECTION !! -->

and you run

	$ ifit -i input.html -o output.html IOS_SECTION OtherSection

then you will include the IOS_SECTION.  If you run

	$ ifit -i input.html -o output.html ANDROID_SECTION OtherSection

you will not include it.

Substitution Values
-------------------

The -s option allows you to read in a JSON file of subduction values.  These are also taken to be
sections that you would want to have turned on.

Example:

	{
		"aa": "[[i am an aa]]",
		"bob": "[[i am a bob]]",
	}

will sububstitue `$$aa$$` for `[[i am an aa]]`.  It will also turn on the section `<!-- !! if aa !! -->`.

Syntax
------

HTML

	<!-- !! if NameA !! -->
	<div id="A"></div>
	<!-- !! end NameA !! -->

or

	<!-- !! if NameB !! 
	<div id="B"></div>
	!! end NameB !! -->

JavaScript

	// !! if NameA !!
	var A = 12;
	// !! end NameA !! 

CSS
	/* !! if NameC !! */
	.classC {
	}
	/* !! end NameC !!  */

or

	/* !! if NameD !!
	.classD {
	}
	!! end NameD !!  */


Please Note
-----------

Tests are in a `Makefile` and run by 

	$ make test1
	$ make test2
	$ make test3

You should see *PASS* at the end of each successful test.

LICENSE
-------

MIT Licnesed -  See LICENSE file.

Author
------

By Philip Schlump.

