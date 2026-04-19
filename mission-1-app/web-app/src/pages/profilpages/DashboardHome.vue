<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > TABLEAU DE BORD</p>
            <h1 class="hero-title1">Bonjour {{prenom}} 👋</h1>
            <p class="classic-text">
                Voici un résumé de votre activité sur UpcycleConnect
            </p>
        </div>
        <button class="btn-main-action">+ Déposer une annonce</button>
    </header>

    <div class="stats-grid">
        <div class="card card--score">
            <p class="tag-score">♻ UPCYCLING SCORE</p>
            <div class="score-value">847 <span>pts</span></div>
            <p class="score-level">Niveau : Éco-Héros 🌿</p>
            <div class="score-footer">
                <div class="mini-stat">
                    <strong>23 kg</strong><br />CO2 évité
                </div>
                <div class="mini-stat">
                    <strong>14</strong><br />Objets recyclés
                </div>
                <div class="mini-stat">
                    <strong>€ 340</strong><br/>Économisé
                </div>
            </div>
        </div>
        <div class="card card--white">
            <div class="card-num">3</div>
            <p class="text-dm">Annonces actives</p>
            <span class="badge badge--green">+1 ce mois</span>
        </div>
        <div class="card card--white">
            <div class="card-num2">2</div>
            <p class="text-dm">Dépôts en attente</p>
            <span class="badge badge--orange">EN COURS</span>
        </div>
    </div>

<div class="section-container">
    <div class="section-header">
        <h2>Mes dernières annonces</h2>
        <div class="header-actions">
            <input
                type="text"
                placeholder="Rechercher..."
                class="search-input"
            />
            <button class="btn-secondary">Tous statuts</button>
            <button class="btn-main-action1">+ Nouvelle annonce</button>
        </div>
    </div>
    <table class="data-table">
        <thead>
            <tr>
                <th>OBJET</th>
                <th>CATÉGORIE</th>
                <th>TYPE</th>
                <th>STATUT</th>
                <th>DATE</th>
                <th>ACTIONS</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>Chaise vintage années 60</td>
                <td>Mobilier</td>
                <td><span class="tag-don">DON</span></td>
                <td><span class="status-valid">✓ VALIDÉE</span></td>
                <td>12 fév. 2026</td>
                <td class="actions-cell">
                    <button class="btn-view">Voir</button>
                    <button class="btn-remove">Retirer</button>
                </td>
            </tr>

            <tr>
                <td>Lot de chutes de tissu lin</td>
                <td>Textile</td>
                <td><span class="tag-vente">VENTE 15€</span></td>
                <td>
                    <span class="status-pending">⌛ EN ATTENTE</span>
                </td>
                <td>18 fév. 2026</td>
                <td class="actions-cell">
                    <button class="btn-view">Voir</button>
                    <button class="btn-modify">Modifier</button>
                </td>
            </tr>
        </tbody>
    </table>
</div>

<div class="section-container">
    <div class="section-header">
        <h2>📅 Mon planning — semaine du 23 fév.</h2>
        <button class="btn-secondary">Vue mensuelle</button>
    </div>
    <div class="planning-row">
        <div class="day-box active">
            <span class="day-num">23</span>
            <div class="event event--green">Atelier bois - 14h</div>
        </div>
        <div class="day-box"><span class="day-num">24</span></div>
        <div class="day-box active">
            <span class="day-num">25</span>
            <div class="event event--orange">Dépôt conteneur 11h</div>
        </div>
        <div class="day-box"><span class="day-num">26</span></div>
    </div>
</div>
    <div class="end-grid">
        <div class="section-container-tips">
            <p class="tag-vente">💡 Conseil du jour</p>
    
            <div class="section-header">
                <h2>Transformer un vieux jean en sac</h2>
            </div>
            <p>
                Apprenez à confectionner un sac tote en 30 minutes avec un
                jean usé. Matériel nécessaire : aiguille, fil, ciseaux.
            </p>
            <button class="btn-view">Lire la suite →</button>
        </div>
        <div class="section-container-tips">
            <p class="tag-don">🔔 Notification</p>
    
            <div class="section-header">
                <h2>Votre dépôt a été récupéré !</h2>
            </div>
            <p>
                La chaise vintage que vous avez déposée le 12 fév. a été
                récupérée par un artisan. +50 points Upcycling Score !
            </p>
            <button class="btn-view">Voir le projet →</button>
        </div>
    </div>
</template>

<script setup>

import {ref} from "vue";

const prenom = ref(localStorage.getItem("userPrenom") || 'Invité');
</script>