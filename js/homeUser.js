	
	var apiUrl = '/users/me/';

	$(document).ready(function(){
		
		//Recuperer les elements
		$.ajax({
			url:	apiUrl + 'interests/',
			type:	'GET',
			success: function(data) {
				if (data) {
					$.each(data, function(index, value){
						$('#interests').append('<li class="interest"><a class="delete" data-id="'+ value.Id+'">x</a>  ' + value.Label + '</li>');
					})
				} else {
					alert(data)
				}
			},
			error: function(xhr,status,error) {
				console.log(error);
			}
		});


		//Ajouter un Element d'interet a la liste
		$('#addinterest').on('submit', function(e){
			e.preventDefault();
			var content = $('#addinterest input[type="text"]').val();
			$.ajax({
				url:			apiUrl + 'interests/',
				type:			'POST',
				data:			JSON.stringify({Label:content}),
				contentType:	'application/json',
				dataType:		'json',
				success:		function(data) {
					if (data.Status == 'ok') {
						$('#addinterest input[type="text"]').val('');
					} else if (data.Status) {
						$('#interests').append('<li class="interest"><a class="delete" data-id="'+ data.Status+'">x</a>  ' + content + '</li>');
						$('#addinterest input[type="text"]').val('');
					} else {
						alert(data);
					}
				}
			})
		});

		//Supprimer un element d'interet a la liste
		$('body').on('click', '.delete', function(e){
			e.preventDefault();
			var self = $(this);
			var id = self.attr('data-id');
			$.ajax({
				url:		apiUrl + 'interests/' + id,
				type:		'DELETE',
				success:	function(data) {
					if (data.Status == 'ok') {
						self.parent('.interest').remove();
					} else {
						alert(data);
					}
				}
			})
		});

	});