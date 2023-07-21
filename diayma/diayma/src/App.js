import React, { useState } from 'react';
import axios from 'axios';
import './liste.css'


function AdminInterface() {
  const [vendors, setVendors] = useState([]);
  const [error, setError] = useState('');

  const handleFetchVendors = () => {
    axios.get('http://192.168.0.73:8080/vendeur/lire/Liste')
      .then(response => {
        setVendors(response.data);
      })
      .catch(error => {
        setError('Erreur lors de la récupération des vendeurs');
        console.error('Erreur :', error);
      });
  };

  return (
    <div className='liste'>
        <img src="../public/LogoDiayma" alt="" />
      <h1>Interface administrateur</h1>
      <button className='bouton' onClick={handleFetchVendors}>Afficher les vendeurs</button>

      {error && <p>{error}</p>}

      <h2>Liste des vendeurs</h2>
      <ul className='list'>
        {vendors.map(vendor => (
          <li key={vendor.id}>
            <p>Nom : {vendor.nom}</p>
            <p>Prénom : {vendor.prenom}</p>
            <p>Login : {vendor.login}</p>
            <p>Téléphone : {vendor.telephone}</p>
            <p>Moyen de paiement : {vendor.moyenPaiement}</p>
            <p>Nom de boutique : {vendor.nomBoutique}</p>
            <p>CIN : {vendor.cin}</p>
            <p>Email : {vendor.email}</p>
            <p>Type de vendeur : {vendor.typeVendeur}</p>
            <p className='sep'>______________________________________________________________________</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default AdminInterface;
