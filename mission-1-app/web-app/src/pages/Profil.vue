<template>
    <section class="sidebar">
        <router-link to="/" class="sidebar__logo">
            Upcycle<span class="text-accent">Connect</span>
        </router-link>

        <div class="hero-divider"></div>

        <router-link to="/profil/informations"class="user-profile-card">
            <div class="user-avatar">ML</div>
            <div class="user-info">
                <span class="user-name">Marie Lambert</span>
                <span class="user-meta">Particulier • Score 847 pts</span>
            </div>
        </router-link>

        <nav class="sidebar__menu">
            <p class="sidebar__category">PRINCIPAL</p>
            <router-link to="/profil" class="sidebar__link">
                <span class="icon-box">🏠</span> Tableau de bord
            </router-link>
            <router-link to="/profil/annonces" class="sidebar__link">
                <span class="icon-box">📢</span> Mes annonces
            </router-link>
            <router-link to="/profil/depots" class="sidebar__link">
                <span class="icon-box">📦</span> Mes dépôts conteneurs
            </router-link>
            <router-link to="/annonces" class="sidebar__link">
                <span class="icon-box">🌿</span> Upcycling Score
            </router-link>

            <p class="sidebar__category">SERVICES</p>
            <router-link to="/formations" class="sidebar__link">
                <span class="icon-box">🎓</span> Formations & Ateliers
            </router-link>
            <router-link to="/communaute" class="sidebar__link">
                <span class="icon-box">🗓️</span> Mon planning
            </router-link>
            <router-link to="/communaute" class="sidebar__link">
                <span class="icon-box">💡</span> Espace Conseils
            </router-link>
            <router-link to="/communaute" class="sidebar__link">
                <span class="icon-box">🏪</span> Catalogue offres
            </router-link>

            <p class="sidebar__category">COMMUNAUTÉ</p>
            <router-link to="/formations" class="sidebar__link">
                <span class="icon-box">💬</span> Forums
            </router-link>
            <router-link to="/formations" class="sidebar__link">
                <span class="icon-box">📅</span> Formations & Ateliers
            </router-link>
        </nav>

        <div class="sidebar__footer">
            <a href="#" class="sidebar__link" @click.prevent="handleLogout">
                <span class="icon-box">🚪</span> Se déconnecter
            </a>
        </div>
    </section>
    <section class="main-content">
        <router-view />
    </section>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const link = ref("TABLEAU DE BORD");

const handleLogout = () => {
  sessionStorage.removeItem("userToken");
  sessionStorage.removeItem("userId");
  
  router.push("/connexion").then(() => {
    window.location.reload();
  });
};
</script>

<style>
:root {
    --text-black: #1a1f1c;
    --text-grey: #5a6660;
    --text-white: #ffffff;
    --text-active-green: #2d7a4f;
    --bg-light: #f7f9f7;
    --sidebar-width: 260px;
}

h2 {
    margin: 10px;
}

.dashboard-layout {
    display: flex;
    min-height: 100vh;
    background-color: #0b0f0e;
}

.sidebar {
    width: var(--sidebar-width);
    height: 100vh;
    background-color: #ffffff;
    padding: 20px 15px;
    position: fixed;
    left: 0;
    top: 0;
    display: flex;
    flex-direction: column;
    border-right: 1px solid #e8ebe9;
    overflow-y: auto;
}

.sidebar::-webkit-scrollbar {
    width: 4px;
}

.sidebar::-webkit-scrollbar-thumb {
    background: #d8e5dc;
    border-radius: 10px;
}

.sidebar__logo {
    font-size: 1.2rem;
    font-weight: 800;
    color: var(--text-black);
    padding-left: 10px; 
    cursor: pointer;
    text-decoration: none;
}

.user-profile-card {
    background-color: #eaf4ed;
    padding: 12px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    text-decoration: none;
}

.user-avatar {
    width: 44px;
    height: 44px;
    background-color: #2d7a4f;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    font-size: 1.1rem;
    border-radius: 12px;
    flex-shrink: 0;
}

.user-info {
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.user-name {
    font-weight: 700;
    color: var(--text-black);
    font-size: 0.9rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.user-meta {
    font-size: 0.75rem;
    color: var(--text-grey);
    white-space: nowrap;
}

.text-accent {
    color: #4ade80;
}

.hero-divider {
    width: 100%;
    height: 1px;
    background: linear-gradient(to right, transparent, #d8e5dc, transparent);
    margin: 10px 0;
}

.sidebar__menu {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
}

.sidebar__category {
    font-size: 0.65rem;
    font-family: "Space Mono", monospace;
    color: #a0ada7;
    letter-spacing: 1px;
    margin: 15px 0 5px 10px;
    text-transform: uppercase;
}

.sidebar__link {
    display: flex;
    align-items: center;
    gap: 10px;
    color: var(--text-grey);
    text-decoration: none;
    padding: 8px 10px;
    border-radius: 10px;
    transition: all 0.2s ease;
    font-size: 0.9rem;
    font-weight: 400;
}

.classic-text {
    display: flex;
    align-items: center;
    gap: 10px;
    color: var(--text-grey);
    text-decoration: none;
    border-radius: 10px;
    transition: all 0.2s ease;
    font-size: 0.9rem;
    font-weight: 400;
}

.icon-box {
    width: 30px;
    height: 30px;
    background-color: var(--bg-light);
    border-radius: 8px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1rem;
}

.sidebar__link:hover {
    background-color: rgba(45, 122, 79, 0.1);
    color: var(--text-active-green);
}

.sidebar__link:hover .icon-box {
    background-color: rgba(45, 122, 79, 0.15);
}

.router-link-exact-active {
    background-color: rgba(45, 122, 79, 0.1);
    color: var(--text-active-green) !important;
}

.router-link-exact-active .icon-box {
    background-color: rgba(45, 122, 79, 0.15);
}

.sidebar__footer {
    padding-top: 15px;
    border-top: 1px solid #f0f2f1;
}

.main-content {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 20px;
    min-height: 100vh;
    background-color: var(--bg-light);
}

.sidebar__category2 {
    font-size: 0.65rem;
    font-family: "Space Mono", monospace;
    color: #a0ada7;
    margin: 0;
    letter-spacing: 1px;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-family: "Syne", sans-serif;
    font-weight: 800;
    font-style: normal;
    line-height: 1.1;
    margin-top: 1.5rem;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
}
.btn-main-action {
    background-color: #2d7a4f;
    color: white;
    border: none;
    padding: 12px 20px;
    border-radius: 8px;
    font-weight: 400;
    cursor: pointer;
}

.header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
}

.search-input,
.btn-secondary {
    background-color: #ffffff;
    border: 1px solid #dcdfdc;
    border-radius: 12px;
    padding: 10px 16px;
    font-size: 14px;
    outline: none;
    color: #333;
    height: 40px;
    box-sizing: border-box;
}

.search-input {
    width: 250px;
}

.search-input::placeholder {
    color: #888;
}

.btn-secondary {
    cursor: pointer;
    background-color: #fcfdfc;
    transition: background 0.2s;
}

.btn-secondary:hover {
    background-color: #f3f5f3;
}

.btn-main-action1 {
    background-color: #2d7a4f;
    color: white;
    border: none;
    height: 40px;
    padding: 0 20px;
    border-radius: 12px;
    font-weight: 500;
    cursor: pointer;
    font-size: 14px;
    display: flex;
    align-items: center;
    transition: background 0.2s;
}

.btn-main-action1:hover {
    background-color: #246340;
}

.stats-grid {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr;
    gap: 20px;
    margin-bottom: 30px;
}

.card {
    background: white;
    border-radius: 16px;
    padding: 24px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card--tips {
    background-color: #2d7a4f;
    color: white;
}

.card--score {
    background-color: #2d7a4f;
    color: white;
}

.card--score p {
    margin-top: 0;
}

.score-value {
    font-size: 3rem;
    font-weight: 800;
    font-family: "Syne", sans-serif;
}

.score-footer {
    display: flex;
    justify-content: space-between;
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
}

.card-num {
    font-family: "Syne", sans-serif;
    font-size: 2.5rem;
    font-weight: 800;
    color: #2d7a4f;
}

.mini-stat {
    font-family: "Syne", sans-serif;
}

.text-dm {
    font-family: "DM Sans", sans-serif;
    font-optical-sizing: auto;
    font-size: 1rem;
    color: var(--text-grey);
}
.card-num2 {
    font-family: "Syne", sans-serif;
    font-size: 2.5rem;
    font-weight: 800;
    color: #e8a030;
}

.end-grid {
    display: grid;
    grid-template-columns: 2fr 2fr;
    gap: 20px;
    margin-bottom: 30px;
}

.section-container {
    background: white;
    border-radius: 16px;
    padding: 20px;
    margin-bottom: 25px;
    box-sizing: border-box;
}

.section-container-tips {
    background: white;
    border-radius: 16px;
    padding: 24px;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    box-sizing: border-box;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.section-container-tips p {
    margin-top: 0;
}

.actions-cell {
    display: flex;
    gap: 8px;
    align-items: center;
    justify-content: flex-start;
}

.tag-don,
.tag-vente {
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: bold;
    display: inline-block;
}

.tag-don {
    background: #fff3e0;
    color: #e8a030;
}

.tag-vente {
    background-color: #e6f4ea;
    color: #1e7e34;
}

.tag-score {
    background-color: rgba(75, 171, 115, 0.15);
    color: #4bab73;
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: bold;
    display: inline-block;
}

.status-valid,
.status-pending {
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: bold;
    display: inline-flex;
    align-items: center;
    gap: 5px;
}

.status-valid {
    background-color: #e6f4ea;
    color: #1e7e34;
}

.status-pending {
    background-color: #fff3e0;
    color: #e8a030;
}

.btn-view,
.btn-modify,
.btn-remove {
    border-radius: 8px;
    padding: 6px 14px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.btn-view {
    background-color: #e6f4ea;
    color: #2d7a4f;
    border: none;
}
.btn-view:hover {
    background-color: #d4ebd9;
}

.btn-modify {
    background-color: #ffffff;
    color: #2d7a4f;
    border: 1px solid #2d7a4f;
}
.btn-modify:hover {
    background-color: #2d7a4f;
    color: #ffffff;
}

.btn-remove {
    background-color: #fceaea;
    color: #c62828;
    border: none;
}
.btn-remove:hover {
    background-color: #f9d7d7;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th {
    text-align: left;
    font-size: 0.75rem;
    color: #a0ada7;
    padding: 10px;
    border-bottom: 1px solid #f0f2f1;
}

.data-table td {
    padding: 15px 10px;
    border-bottom: 1px solid #f0f2f1;
    font-size: 0.9rem;
}

.planning-row {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 10px;
}

.day-box {
    background: #fdfdfd;
    border: 1px solid #f0f2f1;
    border-radius: 12px;
    min-height: 100px;
    padding: 10px;
}

.event {
    font-size: 0.75rem;
    padding: 5px;
    border-radius: 6px;
    margin-top: 10px;
    color: white;
}

.event--green {
    background-color: #2d7a4f;
}
.event--orange {
    background-color: #f39c12;
}

.badge {
    font-size: 0.7rem;
    padding: 4px 8px;
    border-radius: 6px;
    font-weight: bold;
}
.badge--orange {
    background: #fff3e0;
    color: #e8a030;
}
.badge--green {
    background: #e8f5e9;
    color: #2e7d32;
}
</style>
