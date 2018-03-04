![Logo](http://tw.greywool.com/i/jivb7.png)

# Protein Tracker with GO
This is a simple project written in Go Language.
* It takes RNA Nucleic Acid sequence as input. Input should be in **FASTA** format. Means that, It should start with '>' and the next corresponding line should contain nucleic acid sequence. There is a simple portion of Protein database and test case in the Data folder.
* Converts RNA into Amino Acid by translating each codon. It **does** care about STOP and START codons. So watch out for tricky sequences.
* Then searches that AA sequence in given Protein Database.

### Optimization
Because searching matching protein is painfully slow sometimes, I used Go's lovely functionality to speed up things. What it does a parallel search in whole database and stops others once it founds a match. This is not the ideal way of doing this, but for the sake of speed, I implemented it with that way. I might change this to parallel search and reporting findings in the future.

![Test](http://tw.greywool.com/i/OXWdC.png)
## Arguments and Usage
In its most simplistic form, it takes single parameter, the location of the input file.
Like this,
```
go run main.go Data/archaea.5.rna.fna

or

go build
.\GoProteinTracker Data/archaea.5.rna.fna
```
However, for who want to change reference Database for more depth search (Default points Data folder), there is also another parameter which directs the Database's location
```
go run main.go Data/archaea.5.rna.fna Data/ProteinDatabase.txt
```
Moreover, once you point out the Database location, you may want to indicate the index of given RNA sequence(Default is 0) This is important if your input contains
more than one sequence. (In example, there are three). Beware that, you can't use this argument without using Database pointing.  
```
go run main.go Data/archaea.5.rna.fna Data/ProteinDatabase.txt 2
```

### NOTES
This version of protein Database is cropped.(To fit in Github). So you probably can't find your protein in there. Use [this link from WWPDB](ftp://ftp.wwpdb.org/pub/pdb/derived_data) instead and download your fresh database.  
This code still in development.

*Feel free to issue any bugs. Or contact with me if you want to contribute*
