<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > PARAMÈTRES > INFORMATIONS
                </p>
                <h1 class="hero-title1">MON PROFIL</h1>
                <p class="classic-text">
                    Gérez vos informations personnelles et les détails de votre
                    compte utilisateur.
                </p>
                <div v-if="errors.length > 0" 
                style="background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; 
                padding: 10px; border-radius: 8px; margin-bottom: 15px;">
                    <ul style="margin: 0; padding-left: 20px;">
                        <li v-for="(err, index) in errors" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>
                <div v-if="successMsg" class="success-box"
                style="background-color: #E2FEE3; border: 1px solid #44EF44; color: #158F3C; 
                padding: 10px; border-radius: 8px; margin-bottom: 15px;">
                    {{ successMsg }}
                </div>
            </div>
            <button class="btn-main-action" @click="updateProfile">
                Mettre à jour le profil
            </button>
        </header>

        <div class="info-layout">
            <div class="info-main-col">
                <section class="info-section">
                    <h2 class="section-title">Identité Personnelle</h2>
                    <div class="input-grid">
                        <div class="input-field">
                            <label>Prénom</label>
                            <input type="text" v-model="form.prenom" />
                        </div>
                        <div class="input-field">
                            <label>Nom</label>
                            <input type="text" v-model="form.nom" />
                        </div>
                        <div class="input-field">
                            <label>Date de naissance</label>
                            <input type="date" v-model="form.date_naissance" />
                        </div>
                        <div class="input-field">
                            <label>Email (Identifiant)</label>
                            <input
                                type="email"
                                v-model="form.mail"
                                disabled
                                class="disabled-input"
                            />
                        </div>
                    </div>
                </section>

                <section class="info-section">
                    <h2 class="section-title">Adresse et Localisation</h2>
                    <div class="input-field">
                        <label>Adresse complète</label>
                        <input
                            type="text"
                            v-model="form.adresse"
                            placeholder="Ex: 15 rue des Artisans"
                        />
                    </div>
                    <div class="input-grid">
                        <div class="input-field">
                            <label>Ville</label>
                            <input type="text" v-model="form.ville" />
                        </div>
                        <div class="input-field">
                            <label>Code Postal</label>
                            <input type="text" v-model="form.code_postal" />
                        </div>
                    </div>
                </section>
            </div>

            <aside class="info-side-col">
                <section class="info-section status-card">
                    <h2 class="section-title">État du compte</h2>
                    <div class="status-item">
                        <span class="status-label">Rôle actuel</span>
                        <span class="badge-role">{{ form.role }}</span>
                    </div>
                    <div class="status-item">
                        <span class="status-label">Statut</span>
                        <span class="status-active">{{ form.statut }}</span>
                    </div>
                    <div class="status-item">
                        <span class="status-label">Langue ID</span>
                        <span>{{ form.id_langue }}</span>
                    </div>
                    <div class="divider"></div>
                    <div class="registration-info">
                        <p>Membre depuis le :</p>
                        <strong>{{ formatDate(form.date_inscription) }}</strong>
                    </div>
                </section>

                <section class="info-section danger-zone">
                    <h2 class="section-title">Sécurité</h2>
                    <button class="btn-outline" @click="router.push('/profil/password')">Changer le mot de passe</button>
                    <p class="warning-text">
                        Le mot de passe doit être modifié régulièrement.
                    </p>
                </section>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from 'vue-router';

const router = useRouter();

const form = ref({
    prenom: "",
    nom: "",
    mail: "",
    adresse: "",
    ville: "",
    code_postal: "",
    date_naissance: "",
    date_inscription: "",
    role: "",
    statut: "",
    id_langue: 1,
});

const errors = ref([]);
const successMsg = ref("");

onMounted(async () => {
    const id = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");

    if (!id || !token) return;

    try {
        const response = await fetch(`http://localhost:8081/users/${id}`, {
            method: "GET",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
        });

        if (response.ok) {
            const data = await response.json();
            form.value = data;
        } else {
            console.error("Erreur lors de la récupération du profil");
        }
    } catch (error) {
        console.error("Erreur réseau :", error);
    }
});

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};

const updateProfile = async () => {
    const userId = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");

    errors.value = [];
    successMsg.value = "";

    if (!userId || !token) return;

    try {
        const response = await fetch(`http://localhost:8081/users/${userId}`, {
            method: "PUT",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(form.value),
        });

        if (response.ok) {
            successMsg.value = "Profil mis à jour avec succès!";
        } else {
          const data = await response.json();
          errors.value = Array.isArray(data) ? data : [data.message || "Une erreur est survenue"];
        }
    } catch (error) {
        errors.value = ["Le serveur est injoignable pour le moment."];
    }
};
</script>

<style scoped>
.info-view {
    padding: 10px;
}

.info-layout {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 25px;
    margin-top: 20px;
}

.info-section {
    background: white;
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 25px;
    border: 1px solid #e8ebe9;
}

.section-title {
    font-size: 1.1rem;
    font-family: "Syne", sans-serif;
    font-weight: 700;
    margin-bottom: 20px;
    color: #1a1f1c;
}

.input-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
}

.input-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 15px;
}

label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #5a6660;
}

input {
    background-color: #fcfdfc;
    border: 1px solid #dcdfdc;
    border-radius: 10px;
    padding: 12px 15px;
    font-size: 0.95rem;
    outline: none;
    transition: border-color 0.2s;
}

input:focus {
    border-color: #2d7a4f;
}

.disabled-input {
    background-color: #f1f3f2;
    color: #888;
    cursor: not-allowed;
}

.status-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
    font-size: 0.9rem;
}

.status-label {
    color: #5a6660;
}

.badge-role {
    background-color: #eaf4ed;
    color: #2d7a4f;
    padding: 4px 10px;
    border-radius: 8px;
    font-weight: bold;
    font-size: 0.8rem;
}

.status-active {
    color: #1e7e34;
    font-weight: bold;
}

.divider {
    height: 1px;
    background: #e8ebe9;
    margin: 20px 0;
}

.registration-info p {
    font-size: 0.8rem;
    color: #5a6660;
    margin: 0;
}

.registration-info strong {
    font-size: 0.95rem;
}

.btn-outline {
    width: 100%;
    background: white;
    border: 1px solid #dcdfdc;
    padding: 10px;
    border-radius: 10px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
}

.btn-outline:hover {
    background: #f7f9f7;
}

.warning-text {
    font-size: 0.75rem;
    color: #a0ada7;
    margin-top: 10px;
    text-align: center;
}

@media (max-width: 1024px) {
    .info-layout {
        grid-template-columns: 1fr;
    }
    .info-side-col {
        order: -1;
    }
}
</style>
