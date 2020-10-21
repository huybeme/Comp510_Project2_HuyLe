
/*
    Huy Le
	Comp510
	Project 2 part 1
		Take in an input and use input for a search query. Pull data from API and decode.

	Project 2 part 2
		Take data from API and create GUI based window instead printing on console.

	Project 2 part 3
		Retrieve additional data from another API source based on a song selected from part 2.
		Display in GUI a list of similar songs(artist and song names)
*/

notes:
* the list of similar songs compiled in part 3 is a list of clickable widgets that currently only prints a note to console.
    clicking one of these widgets does not work properly. cannot seem to control which element of the list is selected by mouse.
* extra credit: select one of the similar songs and display lyrics
    attempted but not completed! not submitting for credit but to save to complete at a later date.
    so far the commented code displays the lyrics of all songs from the list of similar songs
* when running the program, it is likely to "Client.Timeout exceeded" error when selecting too many songs, especially
    when selecting a song quickly. May need to allow program some time to properly pull data from API