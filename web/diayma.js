// Fonction pour récupérer la liste des vendeurs depuis l'API
function getVendeurs() {
    fetch('http://192.168.0.82:8080/vendeur/lire/Liste')
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

// Fonction pour afficher les vendeurs dans l'interface
function displayVendeurs(vendeurs) {
    const vendeursContainer = document.getElementById('vendeurs');

    if (vendeurs.length === 0) {
        vendeursContainer.innerHTML = '<p>Aucun vendeur trouvé.</p>';
    } else {
        let vendeursHTML = '<ul>';
        vendeurs.forEach(vendeur => {
            vendeursHTML += `<li>${vendeur.nom} - ${vendeur.email} <button class="btnSupprimer" data-id="${vendeur.telephone}">Supprimer</button></li>`;
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
    }
}

// Fonction pour supprimer un vendeur
function supprimerVendeur(vendeurId) {
    fetch('http://192.168.0.82:8080/vendeur/supprime/${vendeurId}', {
        method: 'DELETE',
        mode: 'no-cors'
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

// Ajout du gestionnaire d'événement au bouton "Afficher les vendeurs"
const btnAfficher = document.getElementById('btnAfficher');
btnAfficher.addEventListener('click', () => {
    getVendeurs();
});
