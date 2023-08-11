// Fonction pour récupérer la liste des vendeurs depuis l'API
         // Votre code JavaScript pour récupérer et afficher les vendeurs
         function getVendeurs() {
            const audioElement = document.getElementById('audioElement');
            audioElement.src = "chatbot/audio2.mp3";
            audioElement.play();
            fetch('http://192.168.0.69:8080/vendeur/lire/Liste')
                .then(response => response.json())
                .then(data => {
                    // Appel réussi, traitement des données
                    displayVendeurs(data);
                })
                .catch(error => {
                    // Gestion des erreurs
                    console.error('Erreur lors de la récupération des vendeurs:', error);
                });
        }

        function displayVendeurs(vendeurs) {
            const vendeursContainer = document.getElementById('vendeurs');

            if (vendeurs.length === 0) {
                vendeursContainer.innerHTML = '<p>Aucun vendeur trouvé.</p>';
            } else {
                let vendeursHTML = '<ul>';
                vendeurs.forEach(vendeur => {
                    vendeursHTML += `<li data-nom="${vendeur.nom}" data-email="${vendeur.email}" data-telephone="${vendeur.telephone}" data-prenom="${vendeur.prenom}" data-moyenPaiement="${vendeur.moyenPaiement}" data-typeVendeur="${vendeur.typeVendeur}">${vendeur.nom} - ${vendeur.email} <button class="btnSupprimer" data-id="${vendeur.telephone}">Supprimer</button></li>`;
                });
                vendeursHTML += '</ul>';
                vendeursContainer.innerHTML = vendeursHTML;

                // Ajout des gestionnaires d'événements pour les boutons "Supprimer"
                const btnSupprimerList = document.getElementsByClassName('btnSupprimer');
                for (const btnSupprimer of btnSupprimerList) {
                    btnSupprimer.addEventListener('click', () => {
                        const vendeurId = btnSupprimer.getAttribute('data-id');
                        supprimerVendeur(vendeurId);
                    });
                }

                // Ajout des gestionnaires d'événements pour les noms des vendeurs
                const vendeursList = document.getElementById('vendeurs');
                vendeursList.addEventListener('click', (event) => {
                    if (event.target.tagName === 'LI') {
                        const vendeur = {
                            nom: event.target.getAttribute('data-nom'),
                            email: event.target.getAttribute('data-email'),
                            telephone: event.target.getAttribute('data-telephone'),
                            prenom: event.target.getAttribute('data-prenom')
                        };
                        showPopup(vendeur);
                    }
                });
            }
        }

        function supprimerVendeur(vendeurId) {
            fetch(`http://192.168.0.69:8080/vendeur/supprime/${vendeurId}`, {
                method: 'DELETE'
            })
            .then(response => {
                if (response.ok) {
                    // Vendeur supprimé avec succès, recharger la liste des vendeurs
                    getVendeurs();
                } else {
                    console.error('Erreur lors de la suppression du vendeur:', response.status);
                }
            })
            .catch(error => {
                console.error('Erreur lors de la suppression du vendeur:', error);
            });
        }

        const btnAfficher = document.getElementById('btnAfficher');
        btnAfficher.addEventListener('click', () => {
            getVendeurs();
        });

        function showPopup(vendeur) {
            const popup = document.getElementById('popup');
            const overlay = document.getElementById('overlay');
            const vendeurNom = document.getElementById('vendeurNom');
            const vendeurPrenom = document.getElementById('vendeurPrenom');
            const vendeurEmail = document.getElementById('vendeurEmail');
            const vendeurTelephone = document.getElementById('vendeurTelephone');
            const vendeurMoyenPaiement = document.getElementById('vendeurMoyenPaiement');
            const vendeurTypeVendeur = document.getElementById('vendeurTypeVendeur');


            vendeurNom.textContent = vendeur.nom;
            vendeurPrenom.textContent = vendeur.prenom;
            vendeurEmail.textContent = vendeur.email;
            vendeurTelephone.textContent = vendeur.telephone;
            vendeurMoyenPaiement.textContent = vendeur.moyenPaiement;
            vendeurTypeVendeur.textContent = vendeur.typeVendeur

            popup.style.display = 'block';
            overlay.style.display = 'block';

            // Gestionnaire d'événement pour fermer le popup en cliquant sur le bouton de fermeture
            const closePopup = document.getElementById('closePopup');
            closePopup.addEventListener('click', () => {
                popup.style.display = 'none';
                overlay.style.display = 'none';
            });

            // Gestionnaire d'événement pour fermer le popup lorsque l'overlay est cliqué
            overlay.addEventListener('click', () => {
                popup.style.display = 'none';
                overlay.style.display = 'none';
            });
                        
            //audio
                        const audioElement = document.getElementById('audioElement');
                        audioElement.src = "chatbot/audio3.mp3";
                        audioElement.play();
        }