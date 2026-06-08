<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > FORMATIONS > MODIFIER
                </p>
                <h1 class="hero-title1">MODIFIER UNE FORMATION</h1>
                <p class="classic-text">
                    Mettez à jour les informations de votre session de formation.
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
                <div v-if="isLocked" class="warning-box" style="margin-top: 1rem; margin-bottom: 0;">
                    ⚠️ Cette formation a déjà été validée ou possède déjà des inscrits. Les dates, le type, la capacité et le prix sont verrouillés.
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

        <div
            v-if="isFetching"
            class="loading-state"
            style="margin-top: 2rem; text-align: center"
        >
            Chargement des données de la formation...
        </div>

        <form v-else @submit.prevent="handleSubmit" class="create-annonce-form">
            <div class="split-layout">
                <div class="left-column">
                    <div class="form-card">
                        <h2 class="card-title">1. La Formation</h2>
                        <p class="card-subtitle">
                            Décrivez le contenu et le format de votre session.
                        </p>

                        <div class="form-group">
                            <label>Titre de la formation</label>
                            <input
                                v-model="form.titre"
                                type="text"
                                required
                                placeholder="Ex: Introduction au compostage urbain"
                            />
                        </div>

                        <div class="form-group">
                            <label>Description</label>
                            <textarea
                                v-model="form.description"
                                rows="6"
                                required
                                placeholder="Décrivez le programme, les objectifs, le public visé..."
                            ></textarea>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label>Type de formation</label>
                                <select v-model="form.type" required :disabled="isLocked">
                                    <option value="" disabled>
                                        Sélectionner...
                                    </option>
                                    <option value="Atelier">Atelier</option>
                                    <option value="Cours">Cours</option>
                                    <option value="Webinaire">Webinaire</option>
                                </select>
                            </div>

                            <div class="form-group">
                                <label>Capacité maximale</label>
                                <div class="price-input-wrapper">
                                    <input
                                        v-model.number="form.capacite_max"
                                        type="number"
                                        min="1"
                                        required
                                        placeholder="20"
                                        :disabled="isLocked"
                                    />
                                    <span class="currency-symbol">pers.</span>
                                </div>
                            </div>
                        </div>

                        <div class="form-divider">Dates & Disponibilité</div>

                        <div class="form-row">
                            <div class="form-group">
                                <label>Date de début</label>
                                <input
                                    v-model="form.date_debut"
                                    type="datetime-local"
                                    required
                                    :disabled="isLocked"
                                />
                            </div>
                            <div class="form-group">
                                <label>Date de fin</label>
                                <input
                                    v-model="form.date_fin"
                                    type="datetime-local"
                                    required
                                    :disabled="isLocked"
                                />
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="form-group">
                                <label>Statut</label>
                                <select v-model="form.statut" required>
                                    <option value="" disabled>
                                        Sélectionner...
                                    </option>
                                    <option value="Ouvert">Ouvert</option>
                                    <option value="Complet">Complet</option>
                                    <option value="Termine">Terminé</option>
                                    <option value="Annule">Annulé</option>
                                </select>
                            </div>

                            <div class="form-group">
                                <label>Prix unitaire</label>
                                <div class="price-input-wrapper">
                                    <input
                                        v-model.number="form.prix_unitaire"
                                        type="number"
                                        step="0.01"
                                        min="0"
                                        placeholder="0.00"
                                        :disabled="isLocked"
                                    />
                                    <span class="currency-symbol">€</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card">
                        <h2 class="card-title">2. Lieu de la session</h2>
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
import { ref, computed, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();
const loading = ref(false);
const isFetching = ref(true);
const errors = ref([]);
const successMsg = ref("");

const API_URL = "http://localhost:8081";
const formationId = route.params.id;

const currentUserId = computed(() => {
    const storedId =
        sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});

const isLocked = ref(false);

const form = ref({
    titre: "",
    description: "",
    type: "",
    capacite_max: null,
    date_debut: "",
    date_fin: "",
    statut: "Ouvert",
    prix_unitaire: 0,
    adresse: "",
    ville: "",
    code_postal: "",
});

const formatDateForInput = (isoString) => {
    if (!isoString) return "";
    const date = new Date(isoString);
    const tzOffset = date.getTimezoneOffset() * 60000;
    const localISOTime = new Date(date.getTime() - tzOffset)
        .toISOString()
        .slice(0, 16);
    return localISOTime;
};

const fetchFormationData = async () => {
    const token = sessionStorage.getItem("userToken");

    try {
        const res = await fetch(`${API_URL}/formations/${formationId}`, {
            method: "GET",
            headers: { Authorization: token },
        });

        if (res.ok) {
            const data = await res.json();
            form.value.titre = data.titre || "";
            form.value.description = data.description || "";
            form.value.type = data.type || "";
            form.value.capacite_max = data.capacite_max || null;
            form.value.prix_unitaire = data.prix_unitaire || 0;
            form.value.date_debut = formatDateForInput(data.date_debut);
            form.value.date_fin = formatDateForInput(data.date_fin);
            form.value.adresse = data.adresse || "";
            form.value.code_postal = data.code_postal || "";
            form.value.ville = data.ville || "";
            form.value.statut = data.statut || "Ouvert";
            
            const nbInscrits = data.nb_inscrit || 0;
            isLocked.value = (data.est_valide === "Valide") || (nbInscrits > 0);
            
        } else {
            errors.value.push("Impossible de charger les informations.");
        }
    } catch (err) {
        console.error("Erreur de chargement:", err);
    } finally {
        isFetching.value = false;
    }
};

onMounted(() => {
    fetchFormationData();
});

const validateFrontend = () => {
    errors.value = [];
    if (!form.value.titre || form.value.titre.length < 5)
        errors.value.push("Le titre doit faire au moins 5 caractères.");
    if (!form.value.description || form.value.description.length < 10)
        errors.value.push("La description doit faire au moins 10 caractères.");
    if (!form.value.type)
        errors.value.push("Veuillez sélectionner un type de formation.");
    if (!form.value.capacite_max || form.value.capacite_max < 1)
        errors.value.push("La capacité maximale doit être d'au moins 1.");
    if (!form.value.date_debut)
        errors.value.push("La date de début est obligatoire.");
    if (!form.value.date_fin)
        errors.value.push("La date de fin est obligatoire.");
    if (
        form.value.date_debut &&
        form.value.date_fin &&
        form.value.date_fin <= form.value.date_debut
    )
        errors.value.push(
            "La date de fin doit être postérieure à la date de début.",
        );
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
        id: Number(formationId),
        id_formateur: currentUserId.value,
        titre: form.value.titre,
        description: form.value.description,
        type: form.value.type,
        capacite_max: parseInt(form.value.capacite_max),
        date_debut: new Date(form.value.date_debut).toISOString(),
        date_fin: new Date(form.value.date_fin).toISOString(),
        statut: form.value.statut,
        prix_unitaire: parseFloat(form.value.prix_unitaire) || 0,
        adresse: form.value.adresse || "",
        ville: form.value.ville || "",
        code_postal: form.value.code_postal || "",
    };

    try {
        const response = await fetch(`${API_URL}/formation/${formationId}`, {
            method: "PUT",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            successMsg.value = "Formation modifiée avec succès !";
            setTimeout(() => {
                router.push("/profil/formations");
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

.warning-box {
    background-color: #fffbeb;
    border: 1px solid #f59e0b;
    color: #b45309;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
    font-size: 0.9rem;
    font-weight: 500;
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

input:disabled,
select:disabled {
    background-color: #f5f5f5;
    color: #888;
    cursor: not-allowed;
    border-color: #eee;
}

textarea {
    resize: vertical;
}

.price-input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
}

.price-input-wrapper input {
    width: 100%;
    padding-right: 3rem;
}

.currency-symbol {
    position: absolute;
    right: 1rem;
    font-weight: bold;
    color: #888;
    font-size: 0.9rem;
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