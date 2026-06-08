<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > ÉVÉNEMENTS > MODIFIER</p>
                <h1 class="hero-title1">MODIFIER L'ÉVÉNEMENT</h1>
                <p class="classic-text">
                    Mettez à jour les informations de votre événement.
                </p>

                <div
                    v-if="errors.length > 0"
                    class="error-box"
                    style="margin-top: 1rem; margin-bottom: 0"
                >
                    <ul style="margin: 0; padding-left: 20px">
                        <li v-for="(err, index) in errors" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>

                <div
                    v-if="successMsg"
                    class="success-box"
                    style="margin-top: 1rem; margin-bottom: 0"
                >
                    {{ successMsg }}
                </div>
            </div>

            <button class="btn-secondary" @click="$router.back()">
                🠔 Retour
            </button>
        </header>

        <div v-if="isFetching" class="loading-state" style="margin-top: 2rem; text-align: center">
            Chargement des données de l'événement...
        </div>

        <form v-else @submit.prevent="handleSubmit" class="create-annonce-form">
            <div class="split-layout">
                <div class="left-column">
                    <div class="form-card">
                        <h2 class="card-title">1. L'Événement</h2>
                        <p class="card-subtitle">
                            Décrivez le contenu et le type de votre événement.
                        </p>

                        <div class="form-group">
                            <label>Titre de l'événement</label>
                            <input
                                v-model="form.titre"
                                type="text"
                                required
                                placeholder="Ex: Collecte de déchets électroniques"
                            />
                        </div>

                        <div class="form-group">
                            <label>Description</label>
                            <textarea
                                v-model="form.description"
                                rows="6"
                                required
                                placeholder="Décrivez le programme, les objectifs..."
                            ></textarea>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label>Type d'événement</label>
                                <select v-model="form.type" required>
                                    <option value="" disabled>Sélectionner...</option>
                                    <option value="Atelier">Atelier</option>
                                    <option value="Collecte">Collecte</option>
                                    <option value="Conference">Conférence</option>
                                    <option value="Echange">Échange</option>
                                </select>
                            </div>

                            <div class="form-group">
                                <label>Date de l'événement</label>
                                <input
                                    v-model="form.date_evenement"
                                    type="datetime-local"
                                    required
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card">
                        <h2 class="card-title">2. Lieu de l'événement</h2>
                        <p class="card-subtitle">
                            Renseignez l'adresse exacte du rendez-vous.
                        </p>

                        <div class="form-group">
                            <label>Adresse</label>
                            <input
                                v-model="form.adresse"
                                type="text"
                                placeholder="Ex: 12 rue des Acacias"
                            />
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label>Code postal</label>
                                <input
                                    v-model="form.code_postal"
                                    type="text"
                                    maxlength="5"
                                    placeholder="75001"
                                />
                            </div>
                            <div class="form-group">
                                <label>Ville</label>
                                <input
                                    v-model="form.ville"
                                    type="text"
                                    placeholder="Paris"
                                />
                            </div>
                        </div>
                    </div>

                    <div class="form-actions-card">
                        <button
                            type="button"
                            class="btn-cancel"
                            @click="$router.back()"
                        >
                            Annuler
                        </button>
                        <button
                            type="submit"
                            class="btn-save"
                            :disabled="loading"
                        >
                            {{
                                loading
                                    ? "Enregistrement..."
                                    : "Enregistrer les modifications"
                            }}
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();
const loading = ref(false);
const isFetching = ref(true);
const errors = ref([]);
const successMsg = ref("");

const API_URL = "http://localhost:8081";
const evtId = route.params.id;

const form = ref({
    titre: "",
    description: "",
    type: "",
    date_evenement: "",
    adresse: "",
    ville: "",
    code_postal: "",
});

const formatDateForInput = (isoString) => {
    if (!isoString) return "";
    const date = new Date(isoString);
    const tzOffset = date.getTimezoneOffset() * 60000;
    return new Date(date.getTime() - tzOffset).toISOString().slice(0, 16);
};

const fetchEvtData = async () => {
    const token = sessionStorage.getItem("userToken");
    try {
        const res = await fetch(`${API_URL}/evenements/${evtId}`, {
            headers: { Authorization: token },
        });
        if (res.ok) {
            const data = await res.json();
            form.value = {
                titre: data.titre || "",
                description: data.description || "",
                type: data.type || "",
                date_evenement: formatDateForInput(data.date_evenement),
                adresse: data.adresse || "",
                ville: data.ville || "",
                code_postal: data.code_postal || "",
            };
        } else {
            errors.value.push("Impossible de charger les informations.");
        }
    } catch (err) {
        console.error("Erreur de chargement:", err);
        errors.value.push("Erreur de chargement.");
    } finally {
        isFetching.value = false;
    }
};

const validateFrontend = () => {
    errors.value = [];
    if (!form.value.titre || form.value.titre.length < 5)
        errors.value.push("Le titre doit faire au moins 5 caractères.");
    if (!form.value.description || form.value.description.length < 10)
        errors.value.push("La description doit faire au moins 10 caractères.");
    if (!form.value.type)
        errors.value.push("Veuillez sélectionner un type d'événement.");
    if (!form.value.date_evenement)
        errors.value.push("La date de l'événement est obligatoire.");
    return errors.value.length === 0;
};

const handleSubmit = async () => {
    if (!validateFrontend()) return;

    loading.value = true;
    successMsg.value = "";
    errors.value = [];

    const token = sessionStorage.getItem("userToken");

    try {
        const res = await fetch(`${API_URL}/evenements/${evtId}`, {
            method: "PUT",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                ...form.value,
                date_evenement: new Date(form.value.date_evenement).toISOString(),
            }),
        });

        if (res.ok) {
            successMsg.value = "Événement modifié avec succès !";
            setTimeout(() => router.push("/profil/evenements"), 1500);
        } else {
            const errMsg = await res.text();
            errors.value.push("Erreur serveur : " + errMsg);
        }
    } catch (err) {
        console.error("Erreur:", err);
        errors.value.push("Impossible de joindre le serveur.");
    } finally {
        loading.value = false;
    }
};

onMounted(fetchEvtData);
</script>

<style scoped>
.page-container {
    padding: 20px;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
}

.header-left {
    display: flex;
    flex-direction: column;
}

.sidebar__category2 {
    font-size: 0.65rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 1.5rem 0 0.5rem;
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin: 0;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 500;
}

.error-box {
    background-color: #fee2e2;
    border: 1px solid #ef4444;
    color: #b91c1c;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.success-box {
    background-color: #f0fdf4;
    border: 1px solid #22c55e;
    color: #166534;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.create-annonce-form {
    width: 100%;
    box-sizing: border-box;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    width: 100%;
}

.left-column,
.right-column {
    display: flex;
    flex-direction: column;
}

.form-card {
    background: #ffffff;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-bottom: 2rem;
    width: 100%;
    box-sizing: border-box;
}

.card-title {
    font-size: 1.4rem;
    font-weight: 700;
    margin: 0;
    color: #1a1a1a;
}

.card-subtitle {
    font-size: 0.9rem;
    color: #666;
    margin: 0.2rem 0 0 0;
}

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.form-group label {
    font-size: 0.95rem;
    font-weight: 600;
    color: #333;
}

input,
select,
textarea {
    padding: 0.9rem;
    border: 1px solid #ddd;
    border-radius: 10px;
    font-family: inherit;
    font-size: 0.95rem;
    width: 100%;
    box-sizing: border-box;
}

input:focus,
select:focus,
textarea:focus {
    outline: none;
    border-color: #2d7a4f;
}

textarea {
    resize: vertical;
}

.info-label {
    font-size: 0.85rem;
    color: #2d7a4f;
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.5px;
    display: block;
    margin-bottom: 0.4rem;
}

.address-preview {
    background: #fcfcfc;
    border: 1px solid #eee;
    border-radius: 10px;
    padding: 1rem 1.2rem;
}

.address-preview-text {
    margin: 0;
    font-size: 0.95rem;
    color: #333;
    line-height: 1.8;
}

.form-actions-card {
    margin-top: 0.5rem;
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
}

.btn-save {
    background-color: #2d7a4f;
    color: white;
    padding: 1rem 2.5rem;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-save:hover:not(:disabled) {
    background-color: #246343;
}

.btn-save:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-cancel {
    background: #ffffff;
    border: 1px solid #ccc;
    color: #666;
    padding: 1rem 2.5rem;
    border-radius: 12px;
    cursor: pointer;
    font-weight: 600;
}

.btn-cancel:hover {
    background-color: #f5f5f5;
    border-color: #999;
}

.loading-state {
    font-style: italic;
    color: #8fa396;
}
</style>