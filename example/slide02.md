#Focus on content
Learn 5 markup rules
to make content look great

# #1
Each paragraph is a slide

@append
@zoom/1.2, top left

# #2
Headlines are equal to markdown syntax
# headline 1
## headline 2
### headline 3
#### headline 4
##### headline 5
###### headline 6

# #3
Emphasize *words*
.surrounding them with *asterisks*

# #4
Indent code with 2 spaces
  /* Day of week: Sakamoto's algorithm */
  int dow(int y, int m, int d)
  {
    static int t[] = {0, 3, 2, 5, 0, 3, 5, 1, 4, 6, 2, 4};
    y -= m < 3;
    return (y + y/4 - y/100 + y/400 + t[m-1] + d) % 7;
  }
@code/javascript

# #5
Lines starting with a dot
disable markup such as
headlines, code or blank lines