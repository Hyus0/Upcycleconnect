<template>
    
    <main class="page-container">
        <SiteNavbar
            :is-authenticated="isLoggedIn"
            :user-name="userName"
            user-role="Particulier"
            :user-score="userScore"
        />
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > FORMATIONS > {{ formation?.titre }}
                </p>
                <h1 class="hero-title1">
                    {{ formation?.titre || "Chargement..." }}
                </h1>
            </div>
            <div class="header-actions">
                <button class="btn-secondary" @click="$router.back()">
                    🠔 Retour
                </button>
            </div>
        </header>

        <div v-if="loading" class="loading-state">
            Récupération des données...
        </div>

        <div v-else-if="formation?.id" class="split-layout">
            <div class="left-column">
                <div class="info-card">
                    <div class="card-header-flex">
                        <h2 class="card-title">Présentation de la session</h2>
                        <span
                            :class="[
                                'type-badge',
                                'type-' + formation.type?.toLowerCase(),
                            ]"
                        >
                            {{ formation.type?.toUpperCase() }}
                        </span>
                    </div>

                    <div class="description-section">
                        <label class="info-label"
                            >Description & Programme</label
                        >
                        <div class="description-box">
                            {{ formation.description }}
                        </div>
                    </div>

                    <div class="specs-grid">
                        <div class="spec-item">
                            <label>Date de début</label>
                            <p class="highlight-val">
                                {{ formatDate(formation.date_debut) }}
                            </p>
                        </div>
                        <div class="spec-item">
                            <label>Date de fin</label>
                            <p>{{ formatDate(formation.date_fin) }}</p>
                        </div>
                        <div class="spec-item">
                            <label>Capacité Max</label>
                            <p>{{ formation.capacite_max }} participants</p>
                        </div>
                        <div class="spec-item">
                            <label>Format</label>
                            <p>{{ formation.type }}</p>
                        </div>
                    </div>

                    <div class="location-section">
                        <label>Lieu exact de rendez-vous</label>
                        <p class="address-text">
                            <strong>📍 {{ formation.adresse }}</strong
                            ><br />
                            {{ formation.code_postal }} {{ formation.ville }}
                        </p>
                    </div>
                </div>
            </div>

            <div class="right-column">
                <div class="info-card status-card">
                    <h3>Organisateur de la session</h3>
                    <div class="trainer-preview">
                        <div class="mini-avatar">
                            {{ formation.prenom_formateur?.charAt(0)
                            }}{{ formation.nom_formateur?.charAt(0) }}
                        </div>
                        <div>
                            <p class="trainer-name">
                                {{ formation.prenom_formateur }}
                                {{ formation.nom_formateur }}
                            </p>
                            <button
                                class="link-btn"
                                @click="viewProfile(formation.id_formateur)"
                            >
                            </button>
                        </div>
                    </div>
                </div>

                <div class="info-card dates-card">
                    <h3>Inscriptions</h3>
                    <div class="data-row">
                        <span class="data-label">Validation Admin :</span>
                        <span class="text-success">✓ APPROUVÉ</span>
                    </div>
                    <div class="data-row">
                        <span class="data-label">Places occupées :</span>
                        <span class="status-badge"
                            >{{ formation.nb_inscrit }} /
                            {{ formation.capacite_max }}</span
                        >
                    </div>
                </div>

                <div class="price-card">
                    <label>PRIX DE LA RÉSERVATION</label>
                    <div class="price-value">
                        {{
                            formation.prix_unitaire > 0
                                ? formation.prix_unitaire + "€"
                                : "GRATUIT"
                        }}
                    </div>
                    <p class="price-hint">par personne</p>
                </div>

                <div class="card registration-card">
                    <button
                        @click="handleInscription"
                        class="btn-main-action"
                        :class="{ 'btn-registered': isRegistered }"
                        :disabled="
                            isRegistered ||
                            formation.statut !== 'Ouvert' ||
                            isRegistering
                        "
                    >
                        <span v-if="isRegistered">✓ Déjà inscrit</span>
                        <span v-else>
                            {{
                                formation.statut === "Ouvert"
                                    ? "Réserver ma place"
                                    : "❌ Complet"
                            }}
                        </span>
                    </button>

                    <button
                        v-if="isRegistered"
                        @click="handleQuit"
                        class="btn-quit"
                        :disabled="isLeaving"
                    >
                        Se désister de la session
                    </button>
                </div>
            </div>
        </div>
    </main>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";

import SiteNavbar from "../components/SiteNavbar.vue";

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const isRegistering = ref(false);
const formation = ref(null);
const userScore = ref(0); 

const isRegistered = ref(false);
const isLeaving = ref(false);

const isLoggedIn = computed(() => {
    return !!sessionStorage.getItem("userToken");
});

const userName = computed(() => {
    const prenom = sessionStorage.getItem("userPrenom") || "";
    const nom = sessionStorage.getItem("userNom") || "";
    return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const formatDate = (d) => {
    if (!d || d.startsWith("0001")) return "Non définie";
    return new Date(d).toLocaleString("fr-FR", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    });
};

const fetchDetail = async () => {
    const id = route.params.id;
    
    const userId = sessionStorage.getItem("userId") || 0;

    try {
        const res = await fetch(
            `http://localhost:8081/formations/${id}?user_id=${userId}`,
        );
        if (res.ok) {
            const data = await res.json();
            formation.value = data;
            isRegistered.value = data.is_registered;
        }
    } catch (error) {
        console.error("Erreur fetch :", error);
    } finally {
        loading.value = false;
    }
};

const viewProfile = (id) => {
    router.push({ name: "public-profile", params: { id } });
};

const handleInscription = async () => {
    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");

    if (!token || !userId) {
        alert("Vous devez être connecté pour vous inscrire.");
        return router.push("/connexion");
    }

    isRegistering.value = true;

    try {
        const res = await fetch(
            `http://localhost:8081/api/formations/${formation.value.id}/join`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    id_utilisateur: parseInt(userId),
                }),
            },
        );

        if (res.status === 201) {
            alert("Inscription réussie !");
            isRegistered.value = true;
            fetchDetail();
        } else if (res.status === 409) {
            alert("Désolé, cette formation est déjà complète.");
        } else {
            const errorMsg = await res.text();
            alert("Erreur : " + errorMsg);
        }
    } catch (error) {
        console.error("Erreur lors de l'inscription :", error);
        alert("Impossible de joindre le serveur.");
    } finally {
        isRegistering.value = false;
    }
};

const handleQuit = async () => {
    if (!confirm("Voulez-vous vraiment vous désinscrire de cette formation ?"))
        return;

    const token = sessionStorage.getItem("userToken");
    const userId = sessionStorage.getItem("userId");

    isLeaving.value = true;
    try {
        const res = await fetch(
            `http://localhost:8081/api/formations/${formation.value.id}/quit`,
            {
                method: "POST",
                headers: {
                    Authorization: token,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id_utilisateur: parseInt(userId) }),
            },
        );

        if (res.ok) {
            alert("Vous avez bien été désinscrit.");
            isRegistered.value = false;
            fetchDetail();
        } else {
            alert("Erreur lors de la désinscription.");
        }
    } catch (e) {
        alert("Erreur de connexion.");
    } finally {
        isLeaving.value = false;
    }
};

onMounted(fetchDetail);
</script>

<style scoped>
.page-container {
    min-height: 100vh;
    padding: 20px;
    background: #f7f9f7;
    max-width: 1600px; 
    margin: 0 auto;
}

.content-header {
    padding-top: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 2rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
}

.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    margin: 0.5rem 0 0;
}

.card-header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    border-bottom: 1px solid #f5f5f5;
    padding-bottom: 1rem;
}

.type-badge {
    padding: 6px 14px;
    border-radius: 20px;
    font-weight: 800;
    font-size: 0.75rem;
}

.type-atelier {
    background: #e0f2f1;
    color: #00796b;
}
.type-cours {
    background: #e8eaf6;
    color: #3f51b5;
}
.type-webinaire {
    background: #f3e5f5;
    color: #7b1fa2;
}

.description-section {
    margin-bottom: 2rem;
}

.info-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    text-transform: uppercase;
    color: #999;
    margin-bottom: 8px;
    letter-spacing: 0.5px;
}

.description-box {
    background: #f9f9f9;
    padding: 1.2rem;
    border-radius: 12px;
    color: #444;
    line-height: 1.6;
    font-size: 1rem;
    border: 1px solid #f0f0f0;
}

.specs-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    background: #fcfcfc;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #f5f5f5;
    margin-bottom: 2rem;
}

.spec-item label {
    font-size: 0.7rem;
    color: #aaa;
    text-transform: uppercase;
    margin-bottom: 4px;
    display: block;
}

.spec-item p {
    font-weight: bold;
    color: #333;
    margin: 0;
}
.highlight-val {
    color: #2d6a4f;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    margin-bottom: 3rem;
}

.info-card,
.status-card,
.dates-card,
.price-card {
    background: white;
    padding: 1.5rem;
    border-radius: 16px;
    border: 1px solid #eee;
    margin-bottom: 1.5rem;
}

.trainer-preview {
    display: flex;
    gap: 15px;
    align-items: center;
}
.mini-avatar {
    width: 50px;
    height: 50px;
    background: #2d6a4f;
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
}
.trainer-name {
    font-weight: 800;
    margin: 0;
}
.link-btn {
    background: none;
    border: none;
    color: #2d6a4f;
    text-decoration: underline;
    padding: 0;
    cursor: pointer;
    font-size: 0.8rem;
}

.data-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    font-size: 0.9rem;
    align-items: center;
}
.data-label {
    color: #888;
}
.text-success {
    color: #2d6a4f;
    font-weight: bold;
}
.status-badge {
    background: #f5f5f5;
    padding: 4px 8px;
    border-radius: 4px;
    font-weight: bold;
    font-size: 0.8rem;
}

.price-card {
    text-align: center;
    border-color: #f0f0f0;
}

.price-value {
    font-size: 3.5rem;
    font-weight: 900;
    color: #2d6a4f;
    margin: 10px 0;
}

.price-hint {
    font-size: 0.95rem;
    color: #888;
}

.btn-main-action {
    width: 100%;
    background: #2d6a4f;
    color: white;
    padding: 1.2rem;
    border-radius: 12px;
    font-weight: bold;
    cursor: pointer;
    border: none;
    transition: 0.2s;
}
.btn-main-action:disabled {
    background: #ccc;
    cursor: not-allowed;
}
.btn-main-action:hover:not(:disabled) {
    background: #1b4332;
}

.btn-quit {
    width: 100%;
    background: none;
    border: 1px solid #ff4d4d;
    color: #ff4d4d;
    padding: 10px;
    border-radius: 12px;
    margin-top: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: 0.2s;
}

.btn-quit:hover {
    background: #fff5f5;
    transform: scale(0.98);
}
.sidebar__category2 {
    font-size: 0.75rem;
    color: #999;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}
.btn-secondary {
    background: #f5f5f5;
    color: #666;
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
}

.btn-main-action.btn-registered {
    background-color: #2d6a4f;
    opacity: 0.5;
    cursor: default;
    transform: none;
    box-shadow: none;
}

.btn-modify {
    background: #e8f5e9;
    color: #2d6a4f;
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
    margin-left: 10px;
}
</style>
