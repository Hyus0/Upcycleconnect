<template>
    <header class="content-header">
        <div class="header-left">
            <p class="sidebar__category2">ACCUEIL > ANNONCES > MODIFIER</p>
            <h1 class="hero-title1">MODIFIER L'OBJET</h1>
            <p class="classic-text">Modifiez les détails de votre annonce ci-dessous.</p>
            
            <div v-if="errors.length > 0" class="error-box">
                <ul style="margin: 0; padding-left: 20px;">
                    <li v-for="(err, index) in errors" :key="index">{{ err }}</li>
                </ul>
            </div>
            <div v-if="successMsg" class="success-box">
                {{ successMsg }}
            </div>
        </div>
        <button class="btn-secondary" @click="$router.back()">🠔 Retour</button>
    </header>

    <div class="section-container form-container">
        <div v-if="loading" class="loading-state">Chargement des données...</div>

        <form v-else @submit.prevent="handleUpdate" class="edit-form">
            <div class="form-group">
                <label>Titre de l'annonce</label>
                <input v-model="annonce.titre" type="text" required placeholder="Ex: Chaise en bois vintage" />
            </div>

            <div class="form-group">
                <label>Description détaillée</label>
                <textarea v-model="annonce.description" rows="5" placeholder="Décrivez l'état..."></textarea>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label>Type d'annonce (non modifiable)</label>
                    <select v-model="annonce.type" disabled class="input-disabled">
                        <option value="Vente">Vente</option>
                        <option value="Don">Don</option>
                    </select>
                </div>

                <div class="form-group" v-if="annonce.type === 'Vente'">
                    <label>Prix (€)</label>
                    <input v-model.number="annonce.prix" type="number" step="0.01" />
                </div>

                <div class="form-group">
                    <label>État de l'objet</label>
                    <select v-model="annonce.etat_objet">
                        <option value="Neuf">Neuf</option>
                        <option value="Bon etat">Bon état</option>
                        <option value="Usage">Usagé</option>
                    </select>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label>Ville</label>
                    <input v-model="annonce.ville" type="text" />
                </div>
                <div class="form-group">
                    <label>Code Postal</label>
                    <input v-model="annonce.code_postal" type="text" maxlength="5" />
                </div>
            </div>

            <div class="form-group">
                <label>Adresse (optionnel)</label>
                <input v-model="annonce.adresse" type="text" />
            </div>

            <div class="form-actions">
                <button type="button" class="btn-cancel" @click="$router.back()">Annuler</button>
                <button type="submit" class="btn-save" :disabled="submitting">
                    {{ submitting ? 'Enregistrement...' : 'Enregistrer les modifications' }}
                </button>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const submitting = ref(false);

const errors = ref([]);
const successMsg = ref("");

const annonce = ref({
    titre: '',
    description: '',
    type: 'Vente',
    prix: 0,
    etat_objet: 'Bon etat',
    adresse: '',
    ville: '',
    code_postal: '',
    statut: 'Disponible'
});

onMounted(async () => {
    const idAnnonce = route.params.id;
    const token = localStorage.getItem("userToken");

    try {
        const response = await fetch(`http://localhost:8081/annonces/${idAnnonce}`, {
            headers: { "Authorization": token }
        });

        if (response.ok) {
            const data = await response.json();
            annonce.value = data;
        } else {
            alert("Impossible de récupérer l'annonce.");
            router.push('/annonces');
        }
    } catch (err) {
        console.error("Erreur fetch:", err);
    } finally {
        loading.value = false;
    }
});

const handleUpdate = async () => {
    submitting.value = true;
    const idAnnonce = route.params.id;
    const token = localStorage.getItem("userToken");

    try {
        const response = await fetch(`http://localhost:8081/annonces/${idAnnonce}`, {
            method: "PUT",
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify(annonce.value)
        });

        if (response.ok) {
            successMsg.value = "Votre annonce a été mis à jour !";
        } else {
            const error = await response.json();
            alert("Erreur : " + (errorData.message || "La mise à jour a échoué."));
        }
    } catch (err) {
        console.error("Erreur lors de l'envoi:", err);
        alert("Serveur injoignable.");
    } finally {
        submitting.value = false;
    }
};
</script>

<style scoped>
.edit-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-top: 2rem;
}

.form-row {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.form-group input, 
.form-group select, 
.form-group textarea {
    padding: 0.8rem;
    border: 1px solid #ddd;
    border-radius: 8px;
}

.input-disabled {
    background-color: #f5f5f5;
    color: #888;
    cursor: not-allowed;
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
}

.btn-save {
    background-color: #2e7d32;
    color: white;
    padding: 0.8rem 2rem;
    border: none;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
}

.btn-save:disabled {
    background-color: #a5d6a7;
    cursor: not-allowed;
}

.btn-cancel {
    background: none;
    border: 1px solid #ccc;
    padding: 0.8rem 2rem;
    border-radius: 8px;
    cursor: pointer;
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
</style>