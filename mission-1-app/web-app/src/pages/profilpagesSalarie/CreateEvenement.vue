<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > ÉVÉNEMENTS > CRÉER</p>
                <h1 class="hero-title1">CRÉER UN ÉVÉNEMENT</h1>
                <p class="classic-text">
                    Proposez un événement communautaire autour de l'upcycling.
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

        <form @submit.prevent="handleSubmit" class="create-annonce-form">
            <div class="split-layout">
                <div class="left-column">
                    <div class="form-card">
                        <h2 class="card-title">1. L'Événement</h2>
                        <p class="card-subtitle">
                            Décrivez le contenu et le format de votre rassemblement.
                        </p>

                        <div class="form-group">
                            <label>Titre de l'événement</label>
                            <input
                                v-model="form.titre"
                                type="text"
                                required
                                placeholder="Ex: Grande collecte de palettes"
                            />
                        </div>

                        <div class="form-group">
                            <label>Description</label>
                            <textarea
                                v-model="form.description"
                                rows="6"
                                required
                                placeholder="Décrivez le programme, le but de la rencontre..."
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
                        </div>

                        <div class="form-divider">Date & Heure</div>

                        <div class="form-group">
                            <label>Date de l'événement</label>
                            <input
                                v-model="form.date_evenement"
                                type="datetime-local"
                                required
                            />
                            <div class="form-group">
                                <label>Date de fin</label>
                                <input
                                    v-model="form.date_fin"
                                    type="datetime-local"
                                    required
                                />
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card">
                        <h2 class="card-title">2. Lieu du rendez-vous</h2>
                        <p class="card-subtitle">
                            Renseignez l'adresse exacte pour les participants.
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
                            {{ loading ? "Envoi en cours..." : "Soumettre l'événement" }}
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const loading = ref(false);
const errors = ref([]);
const successMsg = ref("");

const API_URL = "/go";

const currentUserId = computed(() => {
    const storedId =
        sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});

const form = ref({
    titre: "",
    description: "",
    type: "",
    date_evenement: "",
    date_fin: "",
    adresse: "",
    ville: "",
    code_postal: "",
});

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
    if (!form.value.date_fin)
      errors.value.push("La date de fin est obligatoire.");
    if (
        new Date(form.value.date_fin) <=
        new Date(form.value.date_evenement)
    ) {
        errors.value.push(
            "La date de fin doit être postérieure à la date de début."
        );
    }
    return errors.value.length === 0;
};

const handleSubmit = async () => {
    if (currentUserId.value === 0) {
        errors.value = ["Erreur de session. Reconnectez-vous."];
        return;
    }

    if (!validateFrontend()) return;

    loading.value = true;
    successMsg.value = "";
    errors.value = [];

    const token = sessionStorage.getItem("userToken");

    const payload = {
        id_createur: currentUserId.value,
        titre: form.value.titre,
        description: form.value.description,
        type: form.value.type,
        date_evenement: new Date(form.value.date_evenement).toISOString(),
        date_fin: new Date(form.value.date_fin).toISOString(),
        adresse: form.value.adresse || "",
        ville: form.value.ville || "",
        code_postal: form.value.code_postal || "",
    };

    try {
        const response = await fetch(`${API_URL}/evenements`, {
            method: "POST",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            successMsg.value = "Événement créé avec succès !";
            setTimeout(() => {
                router.push("/profil/evenements");
            }, 1500);
        } else {
            const errMsg = await response.text();
            errors.value.push("Erreur serveur : " + errMsg);
        }
    } catch (err) {
        console.error("Erreur:", err);
        errors.value.push("Impossible de joindre le serveur.");
    } finally {
        loading.value = false;
    }
};
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

.form-divider {
    margin-top: 0.5rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e9f5ed;
    font-weight: 800;
    color: #2d7a4f;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 1px;
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
</style>