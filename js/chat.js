

$('#textin').val("");
// take what's the textbox and send it off
$('#send').click( function(event) {
	sock.send($('#textin').val());
	$('#textin').val("");
});