<template>
    <div class="layout-wrapper public-dashboard">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    {{ t.Home || "Accueil" }} >
                    {{
                        t.DashboardPro
                            ? t.DashboardPro.toUpperCase()
                            : "ESPACE PROFESSIONNEL"
                    }}
                </p>
                <h1 class="hero-title1">
                    {{ t.Hello || "Bonjour" }} {{ prenom }} 👋
                </h1>
                <p class="classic-text">
                    {{
                        t.SummaryActivityPro ||
                        "Voici le tableau de bord avancé de votre activité sur UpcycleConnect."
                    }}
                </p>
            </div>
            <div class="header-actions-group">
                <span class="premium-pill">
                    <i class="ti ti-crown" aria-hidden="true"></i>
                    Compte Premium
                </span>
                <router-link to="/profil/createAnnonce" class="btn-main-action">{{
                    t.DeposeAnnonce || "Déposer une annonce"
                }}</router-link>
            </div>
        </header>

        <div class="stats-grid">
            <div class="card card--score">
                <p class="tag-score">UPCYCLING SCORE PRO</p>
                <div class="score-value">
                    {{ stats.total_points || 0 }} <span>pts</span>
                </div>
                <p class="score-level">
                    {{ t.Niveau || "Niveau" }} : {{ stats.niveau || "Partenaire" }}
                </p>
                <div class="score-footer">
                    <div class="mini-stat">
                        <strong>{{ stats.co2_total_evite_kg || 0 }} kg</strong
                        ><br />{{ t.Co2 || "CO2" }}
                    </div>
                    <div class="mini-stat">
                        <strong>{{ stats.nb_objets_recycles || 0 }}</strong
                        ><br />{{ t.Objets || "Objets" }}
                    </div>
                    <div class="mini-stat">
                        <strong
                            >{{ stats.ressources_economisees || 0 }} €</strong
                        ><br />{{ t.Economise || "Économisés" }}
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-num">{{ projets.length }}</div>
                <p class="text-dm">Projets mis en avant</p>
                <span class="badge badge--green">Visibilité boostée</span>
            </div>

            <div class="card">
                <div class="card-num" style="color: #d32f2f">
                    {{ recuperationsEnAttente }}
                </div>
                <p class="text-dm">Récupérations en attente</p>
                <span class="badge badge--orange">{{
                    recuperationsEnAttente > 0 ? "À récupérer" : "Aucun retrait"
                }}</span>
            </div>

            <div class="card">
                <div class="card-num" style="color: #2d7a4f">
                    {{ abonnement.statut === 'Actif' ? '✓' : '—' }}
                </div>
                <p class="text-dm">Abonnement {{ abonnement.nom || 'Premium' }}</p>
                <span class="badge" :class="abonnement.statut === 'Actif' ? 'badge--green' : 'badge--orange'">
                    {{ abonnement.statut === 'Actif' ? "Actif jusqu'au " + formatDate(abonnement.date_fin) : "Renouvellement requis" }}
                </span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>Analyse d'impact écologique détaillée</h2>
                    <p class="classic-text">
                        Suivi avancé réservé aux comptes professionnels.
                    </p>
                </div>
                <span class="badge badge--green">
                    <i class="ti ti-leaf" aria-hidden="true"></i> Données 12 derniers mois
                </span>
            </div>

            <div class="eco-metrics-grid">
                <div class="eco-metric-card">
                    <span class="eco-metric-label">CO2 total évité</span>
                    <span class="eco-metric-value">{{ ecoStats.co2_total || 0 }} kg</span>
                    <span class="eco-metric-trend" :class="ecoStats.co2_trend >= 0 ? 'trend-up' : 'trend-down'">
                        <i :class="ecoStats.co2_trend >= 0 ? 'ti ti-trending-up' : 'ti ti-trending-down'" aria-hidden="true"></i>
                        {{ Math.abs(ecoStats.co2_trend || 0) }}% vs mois précédent
                    </span>
                </div>
                <div class="eco-metric-card">
                    <span class="eco-metric-label">Eau économisée</span>
                    <span class="eco-metric-value">{{ ecoStats.eau_economisee || 0 }} L</span>
                    <span class="eco-metric-trend trend-up">
                        <i class="ti ti-droplet" aria-hidden="true"></i>
                        Estimation cumulée
                    </span>
                </div>
                <div class="eco-metric-card">
                    <span class="eco-metric-label">Matériaux valorisés</span>
                    <span class="eco-metric-value">{{ ecoStats.materiaux_valorises || 0 }}</span>
                    <span class="eco-metric-trend trend-up">
                        <i class="ti ti-recycle" aria-hidden="true"></i>
                        Objets transformés
                    </span>
                </div>
                <div class="eco-metric-card">
                    <span class="eco-metric-label">Score impact moyen</span>
                    <span class="eco-metric-value">{{ ecoStats.score_impact_moyen || 0 }} / 100</span>
                    <span class="eco-metric-trend trend-up">
                        <i class="ti ti-chart-bar" aria-hidden="true"></i>
                        Sur vos projets publiés
                    </span>
                </div>
            </div>

            <div class="chart-wrapper">
                <div style="position: relative; width: 100%; height: 260px;">
                    <canvas
                        id="co2Chart"
                        role="img"
                        aria-label="Graphique en barres du CO2 évité par mois sur les 6 derniers mois"
                    >Évolution du CO2 évité sur 6 mois.</canvas>
                </div>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Statistiques des matériaux disponibles</h2>
                <router-link
                    to="/annonces"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir toutes les annonces</router-link
                >
            </div>

            <div class="materials-grid">
                <div class="chart-wrapper" style="flex: 1;">
                    <div style="position: relative; width: 100%; height: 240px;">
                        <canvas
                            id="materiauxChart"
                            role="img"
                            aria-label="Graphique en anneau de la répartition des matériaux disponibles par type"
                        >Répartition des matériaux disponibles par type.</canvas>
                    </div>
                    <div class="chart-legend">
                        <span
                            v-for="(item, i) in materiauxStats"
                            :key="i"
                            class="legend-item"
                        >
                            <span class="legend-dot" :style="{ background: materiauxColors[i % materiauxColors.length] }"></span>
                            {{ item.type_materiau }} ({{ item.count }})
                        </span>
                    </div>
                </div>

                <div class="materials-table-wrapper">
                    <table class="data-table">
                        <thead>
                            <tr>
                                <th>Matériau</th>
                                <th>Disponibles</th>
                                <th>Zone</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, i) in materiauxStats.slice(0, 5)" :key="i">
                                <td><strong>{{ item.type_materiau }}</strong></td>
                                <td>
                                    <span class="material-tag">{{ item.count }} dispo.</span>
                                </td>
                                <td class="table-subtext">{{ item.zone || "Multi-zones" }}</td>
                            </tr>
                            <tr v-if="materiauxStats.length === 0">
                                <td colspan="3" class="text-center text-grey py-4">
                                    Aucune donnée matériau disponible.
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>Matériaux recherchés</h2>
                    <p class="classic-text">
                        Indiquez les matériaux qui vous intéressent (séparés par des virgules).
                        Ils servent à calculer vos alertes prioritaires de collecte.
                    </p>
                </div>
            </div>

            <form class="materiaux-form" @submit.prevent="saveMateriauxRecherches">
                <input
                    v-model="materiauxRecherchesInput"
                    type="text"
                    class="materiaux-input"
                    placeholder="Ex: bois, métal, palette, tissu..."
                    maxlength="255"
                />
                <button
                    type="submit"
                    class="btn-main-action"
                    :disabled="savingMateriaux"
                >
                    {{ savingMateriaux ? "Enregistrement..." : "Enregistrer" }}
                </button>
            </form>

            <div v-if="materiauxSavedMessage" class="materiaux-saved-msg">
                ✓ {{ materiauxSavedMessage }}
            </div>

            <div v-if="materiauxKeywordsPreview.length > 0" class="materiaux-keywords-preview">
                <span
                    v-for="(kw, i) in materiauxKeywordsPreview"
                    :key="i"
                    class="material-tag"
                >
                    {{ kw }}
                </span>
            </div>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>Alertes priorisées pour la collecte</h2>
                    <p class="classic-text">
                        Basées sur vos matériaux recherchés et votre localisation.
                    </p>
                </div>
                <span class="badge badge--orange">
                    <i class="ti ti-bell" aria-hidden="true"></i> {{ alertesPrioritaires.length }} active(s)
                </span>
            </div>

            <div v-if="materiauxRecherchesInput.trim() === ''" class="state-card mb-2">
                Renseignez vos matériaux recherchés ci-dessus pour affiner vos alertes prioritaires.
            </div>

            <table class="data-table">
                <thead>
                    <tr>
                        <th>Annonce</th>
                        <th>Type</th>
                        <th>Prix</th>
                        <th>Lieu</th>
                        <th>Date de création</th>
                        <th class="text-right">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="alerte in alertesPrioritaires.slice(0, 4)" :key="alerte.id">
                        <td>
                            <strong>{{ alerte.titre || "Objet" }}</strong><br />
                            <span class="material-tag" style="margin-top: 4px; display: inline-block;">
                                {{ alerte.type_materiau || "Non précisé" }}
                            </span>
                        </td>
                        <td>
                            <span class="badge" :class="alerte.type === 'Don' || alerte.prix == 0 ? 'badge--green' : 'badge--orange'">
                                {{ alerte.type === 'Don' || alerte.prix == 0 ? 'Don' : 'Vente' }}
                            </span>
                        </td>
                        <td>
                            <span class="badge" :class="alerte.type === 'Don' || alerte.prix == 0 ? 'badge--green' : 'badge--orange'">
                                {{ alerte.type === 'Don' || alerte.prix == 0 ? 'Gratuit' : alerte.prix + ' €' }}
                            </span>
                        </td>
                        <td>
                            <span class="status-logistique">{{ alerte.ville || "Non assigné" }}</span>
                        </td>
                        <td>
                            {{ formatDate(alerte.date_creation) }}
                        </td>
                        <td class="text-right">
                            <button
                                class="btn-view"
                                @click="$router.push(`/annonce/${alerte.id}`)"
                            >
                                Voir
                            </button>
                        </td>
                    </tr>
                    <tr v-if="alertesPrioritaires.length === 0">
                        <td colspan="5" class="text-center text-grey py-4">
                            Aucune alerte prioritaire pour le moment.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Récupérations en attente</h2>
                <router-link
                    to="/profil/recuperations"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir tout</router-link
                >
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Objet</th>
                        <th>Conteneur / Site</th>
                        <th>Code de retrait</th>
                        <th>Statut</th>
                        <th class="text-right">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="achat in achats.slice(0, 4)" :key="achat.id">
                        <td>
                            <strong>{{ achat.titre_annonce || "Objet" }}</strong
                            ><br />
                            <small class="table-subtext">{{
                                formatDate(achat.date_transaction)
                            }}</small>
                        </td>
                        <td>
                            {{ achat.nom_site || "En attente d'assignation" }}
                        </td>
                        <td>
                            <strong style="letter-spacing: 2px">{{
                                achat.code_retrait || "---"
                            }}</strong>
                        </td>
                        <td>
                            <span
                                class="badge"
                                :class="
                                    achat.statut === 'Récupéré'
                                        ? 'badge--green'
                                        : 'status-logistique'
                                "
                            >
                                {{ achat.statut || "À récupérer" }}
                            </span>
                        </td>
                        <td class="text-right">
                            <button
                                v-if="achat.statut !== 'Récupéré'"
                                class="btn-view"
                                @click="$router.push('/profil/depots')"
                            >
                                Suivre
                            </button>
                            <span v-else class="text-sm text-gray-500"
                                >Terminé</span
                            >
                        </td>
                    </tr>
                    <tr v-if="achats.length === 0">
                        <td colspan="5" class="text-center text-grey py-4">
                            Vous n'avez aucune récupération en attente.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <h2>Mes Projets</h2>
                <router-link
                    to="/projets"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir le catalogue</router-link
                >
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th style="width: 40%">Projet</th>
                        <th>IMPACT CO2</th>
                        <th>ENGAGEMENT</th>
                        <th>STATUT</th>
                        <th>DATE</th>
                        <th>ACTIONS</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="projet in projets.slice(0, 4)" :key="projet.id">
                        <td>
                            <strong>{{ projet.titre }}</strong
                            ><br />
                            <small class="text-truncate">{{
                                projet.description_courte ||
                                "Aucune description"
                            }}</small>
                        </td>
                        <td>
                            <span class="material-tag">
                                {{ projet.co2_evite_kg || 0 }} kg
                            </span>
                        </td>
                        <td>
                            <span class="stats-text"
                                >{{ projet.nb_vues || 0 }} vues,
                                {{ projet.nb_likes || 0 }} likes</span
                            ><br />
                        </td>
                        <td>
                            <span
                                :class="
                                    projet.visible_public
                                        ? 'status-valid'
                                        : 'status-neutral'
                                "
                            >
                                {{ projet.visible_public ? "PUBLIC" : "PRIVÉ" }}
                            </span>
                        </td>
                        <td>{{ formatDate(projet.date_creation) }}</td>
                        <td class="actions-cell">
                            <button
                                class="btn-view"
                                type="button"
                                @click="$router.push(`/projets/${projet.id}`)"
                            >
                                {{ t.Voir || "Voir" }}
                            </button>
                            <button
                                class="btn-modify"
                                type="button"
                                @click="
                                    $router.push(
                                        `/profil/modifyProjet/${projet.id}`,
                                    )
                                "
                            >
                                {{ t.Modifier || "Modifier" }}
                            </button>
                            <button
                                class="btn-remove"
                                type="button"
                                @click="removeProjet(projet.id)"
                            >
                                {{ t.Supprimer || "Retirer" }}
                            </button>
                        </td>
                    </tr>
                    <tr v-if="projets.length === 0">
                        <td colspan="6" class="text-center text-grey py-4">
                            Vous n'avez créé aucun projet.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>Facturation & abonnement</h2>
                    <p class="classic-text">
                        Gérez votre contrat d'abonnement Premium et vos factures.
                    </p>
                </div>
                <router-link
                    to="/profil/facturation"
                    class="btn-secondary"
                    style="text-decoration: none"
                    >Voir tout</router-link
                >
            </div>
            <table class="data-table">
                <thead>
                    <tr>
                        <th>Référence</th>
                        <th>Commande</th>
                        <th>Montant</th>
                        <th>Statut</th>
                        <th>Date</th>
                        <th class="text-right">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="facture in factures.slice(0, 4)" :key="facture.id">
                        <td><strong>{{ facture.numero_facture || "FAC-" + facture.id }}</strong></td>
                        <td>Commande #{{ facture.commande_id || facture.id }}</td>
                        <td><strong>{{ formatPrice(facture.montant_total || facture.montant) }}</strong></td>
                        <td>
                            <span
                                class="badge"
                                :class="(facture.statut_paiement === 'Payé' || facture.statut === 'Payée') ? 'badge--green' : 'badge--orange'"
                            >
                                {{ facture.statut_paiement || facture.statut || "Payé" }}
                            </span>
                        </td>
                        <td>{{ formatDate(facture.date_transaction || facture.date_emission) }}</td>
                        <td class="actions-cell" style="justify-content: flex-end;">
                            <button class="btn-view" type="button" @click="downloadFacture(facture.id)" title="Télécharger">
                                Télécharger
                            </button>
                            <button class="btn-secondary" type="button" @click="sendFacture(facture.id)" title="Envoyer par mail">
                                Recevoir par mail
                            </button>
                        </td>
                    </tr>
                    <tr v-if="factures.length === 0">
                        <td colspan="6" class="text-center text-grey py-4">
                            Aucune facture pour le moment.
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <div class="section-container">
            <div class="section-header">
                <div>
                    <h2>
                        {{ t.MySchedule || "Mon Planning" }} ({{
                            planningWeekLabel
                        }})
                    </h2>
                    <p class="classic-text">
                        {{
                            t.Schedule ||
                            "Vos prochains événements et formations."
                        }}
                    </p>
                </div>
                <button
                    class="btn-secondary"
                    @click="$router.push('/profil/planning')"
                >
                    {{ t.VueMensuelle || "Vue Mensuelle" }}
                </button>
            </div>

            <div class="planning-week">
                <div
                    class="planning-day"
                    v-for="day in compactWeekDays"
                    :key="day.key"
                    :class="{ 'is-today': day.isToday }"
                >
                    <div class="day-label">
                        {{ day.weekLabel }} <strong>{{ day.dayLabel }}</strong>
                    </div>
                    <div class="entries">
                        <div
                            v-for="entry in day.entries.slice(0, 2)"
                            :key="entry.id"
                            class="entry"
                            :class="`entry--${entry.kind}`"
                        >
                            <strong>{{ entry.title }}</strong>
                            <span>{{ entry.timeLabel }}</span>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="calendarEntries.length === 0" class="state-card mt-4">
                {{ t.PlanningInfo || "Rien de prévu." }}
            </div>
        </div>

        <div class="end-grid">
            <div class="section-container">
                <p class="badge badge--green mb-2">
                    💡 {{ t.Tips || "Conseil du jour" }}
                </p>
                <h2 v-if="tipDuJour">{{ tipDuJour.titre }}</h2>
                <h2 v-else>{{ t.WaitingTips || "En attente de conseil" }}</h2>
                <p class="mt-2">
                    {{
                        tipDuJour
                            ? tipDuJour.description
                            : t.TipsText || "Revenez plus tard."
                    }}
                </p>
                <router-link
                    v-if="tipDuJour"
                    :to="`/conseil/${tipDuJour.id}`"
                    class="btn-text-green mt-4"
                    style="text-decoration: none"
                >
                    {{ t.ReadNext || "Lire la suite" }} →
                </router-link>
            </div>

            <div class="section-container">
                <p class="badge badge--orange mb-2">
                    🔔 {{ t.Notification || "Notification" }}
                </p>
                <h2 v-if="latestNotification">
                    {{ latestNotification.titre }}
                </h2>
                <h2 v-else>
                    {{ t.NothingNotification || "Aucune notification" }}
                </h2>
                <p class="mt-2">
                    {{
                        latestNotification
                            ? latestNotification.message
                            : t.NotificationText || "Rien de nouveau."
                    }}
                </p>
                <button
                    v-if="latestNotification"
                    class="btn-text-green mt-4"
                    @click="handleNotificationClick(latestNotification)"
                >
                    {{ t.SeeAction || "Voir l'action" }} →
                </button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed, onMounted, ref, nextTick } from "vue";
import { useRouter } from "vue-router";
import { useTraduction } from "../../composables/useTraduction";

const router = useRouter();
const API_URL = "/go";
const { t } = useTraduction();

const prenom = ref(sessionStorage.getItem("userPrenom") || "Invité");
const achats = ref([]);
const projets = ref([]);
const calendarEntries = ref([]);
const tipDuJour = ref(null);
const latestNotification = ref(null);
const alertesPrioritaires = ref([]);
const materiauxStats = ref([]);
const factures = ref([]);
const abonnement = ref({ nom: "Premium", statut: "Inactif", date_fin: null });
const ecoStats = ref({
    co2_total: 0,
    co2_trend: 0,
    eau_economisee: 0,
    materiaux_valorises: 0,
    score_impact_moyen: 0,
    co2_par_mois: [],
});

const materiauxRecherchesInput = ref("");
const savingMateriaux = ref(false);
const materiauxSavedMessage = ref("");

const materiauxKeywordsPreview = computed(() =>
    materiauxRecherchesInput.value
        .split(",")
        .map((k) => k.trim())
        .filter(Boolean),
);

const materiauxColors = [
    "#2d7a4f", "#f1a321", "#378ADD", "#D85A30", "#7F77DD", "#D4537E",
];

const stats = ref({
    total_points: 0,
    niveau: "...",
    co2_total_evite_kg: 0,
    nb_objets_recycles: 0,
    ressources_economisees: 0,
});

const recuperationsEnAttente = computed(
    () => achats.value.filter((a) => a.statut !== "Récupéré").length,
);

const getLocalISODate = (date) => {
    const d = new Date(date);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
};

const currentWeekStart = computed(() => {
    const now = new Date();
    now.setHours(0, 0, 0, 0);
    now.setDate(now.getDate() - ((now.getDay() + 6) % 7));
    return now;
});

const planningWeekLabel = computed(() =>
    currentWeekStart.value.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
    }),
);

const weekDays = computed(() => [
    t.value.Lundi || "LUN",
    t.value.Mardi || "MAR",
    t.value.Mercredi || "MER",
    t.value.Jeudi || "JEU",
    t.value.Vendredi || "VEN",
    t.value.Samedi || "SAM",
    t.value.Dimanche || "DIM",
]);

const compactWeekDays = computed(() =>
    Array.from({ length: 7 }, (_, i) => {
        const date = new Date(currentWeekStart.value);
        date.setDate(date.getDate() + i);
        const iso = getLocalISODate(date);
        return {
            key: iso,
            weekLabel: weekDays.value[i].toUpperCase(),
            dayLabel: date.toLocaleDateString("fr-FR", { day: "numeric" }),
            isToday: iso === getLocalISODate(new Date()),
            entries: calendarEntries.value.filter((e) => e.date === iso),
        };
    }),
);

const formatDate = (val) => {
    if (!val) return "---";
    const d = new Date(val);
    return isNaN(d)
        ? "---"
        : new Intl.DateTimeFormat("fr-FR", {
              day: "2-digit",
              month: "short",
              year: "numeric",
          }).format(d);
};

const handleNotificationClick = async (notif) => {
    if (!notif) return;
    const token = sessionStorage.getItem("userToken") || "";
    const userId = sessionStorage.getItem("userId");
    try {
        await fetch(`${API_URL}/notifications/${notif.id}/read`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ user_id: parseInt(userId, 10) }),
        });
    } catch (e) {
        console.error(e);
    }
    router.push(
        notif.titre.toLowerCase().includes("disponible") ||
            notif.type === "Alerte"
            ? "/profil/depots"
            : "",
    );
};

const removeProjet = async (id) => {
    if (
        !confirm(
            t.value.ConfirmDeleteProjet ||
                "Voulez-vous vraiment supprimer ce projet ?",
        )
    )
        return;
    try {
        const token = sessionStorage.getItem("userToken") || "";
        const res = await fetch(`${API_URL}/projets/${id}`, {
            method: "DELETE",
            headers: { Authorization: `Bearer ${token}` },
        });
        if (!res.ok) throw new Error("Erreur suppression projet");
        projets.value = projets.value.filter((p) => p.id !== id);
    } catch (e) {
        console.error("Erreur suppression:", e);
        alert(
            t.value.ErrorDeleteProjet ||
                "Une erreur est survenue lors de la suppression.",
        );
    }
};

const userId = () => parseInt(sessionStorage.getItem("userId"), 10);
const authHeaders = () => ({
    Authorization: `Bearer ${sessionStorage.getItem("userToken") || ""}`,
});

async function loadMateriauxRecherches() {
    const id = userId();
    if (isNaN(id)) return;
    try {
        const res = await fetch(`${API_URL}/users/${id}/materiaux-recherches`, {
            headers: authHeaders(),
        });
        if (res.ok) {
            const data = await res.json();
            materiauxRecherchesInput.value = data.materiaux_recherches || "";
        }
    } catch (e) {
        console.error("Erreur loadMateriauxRecherches:", e);
    }
}

async function saveMateriauxRecherches() {
    const id = userId();
    if (isNaN(id)) return;
    savingMateriaux.value = true;
    materiauxSavedMessage.value = "";
    try {
        const res = await fetch(`${API_URL}/users/${id}/materiaux-recherches`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                ...authHeaders(),
            },
            body: JSON.stringify({
                materiaux_recherches: materiauxRecherchesInput.value.trim(),
            }),
        });
        if (!res.ok) throw new Error("Erreur lors de l'enregistrement");

        materiauxSavedMessage.value = "Matériaux recherchés mis à jour.";
        await loadAlertesPrioritaires();

        setTimeout(() => {
            materiauxSavedMessage.value = "";
        }, 3000);
    } catch (e) {
        console.error("Erreur saveMateriauxRecherches:", e);
        alert("Impossible d'enregistrer vos matériaux recherchés.");
    } finally {
        savingMateriaux.value = false;
    }
}

async function loadAlertesPrioritaires() {
    const id = userId();
    if (isNaN(id)) return;
    try {
        const res = await fetch(`${API_URL}/users/${id}/alertes-prioritaires`, {
            headers: authHeaders(),
        });
        if (res.ok) {
            const data = await res.json();
            alertesPrioritaires.value = Array.isArray(data) ? data : [];
        } else {
            console.error("L'API a renvoyé une erreur :", res.status);
        }
    } catch (e) {
        console.error("Erreur loadAlertesPrioritaires:", e);
    }
}

const formatPrice = (value) => new Intl.NumberFormat("fr-FR", {
    style: "currency",
    currency: "EUR"
}).format(Number(value) || 0);

const downloadFacture = (factureId) => {
    const id = userId();
    if (isNaN(id)) return;
    window.open(`${API_URL}/users/${id}/factures/${factureId}/download`, "_blank");
};

const sendFacture = async (factureId) => {
    const id = userId();
    if (isNaN(id)) return;
    try {
        const res = await fetch(`${API_URL}/users/${id}/factures/${factureId}/send`, {
            method: "POST",
            headers: authHeaders()
        });
        const payload = res.ok ? await res.json() : { message: "Erreur lors de l'envoi." };
        alert(payload.message);
    } catch (error) {
        alert("Impossible de contacter le serveur pour envoyer la facture.");
    }
};

function renderCharts() {
    if (typeof Chart === "undefined") return;

    const co2Canvas = document.getElementById("co2Chart");
    if (co2Canvas) {
        const months = (ecoStats.value.co2_par_mois || []).map((m) => m.mois);
        const values = (ecoStats.value.co2_par_mois || []).map((m) => m.valeur);
        new Chart(co2Canvas, {
            type: "bar",
            data: {
                labels: months.length ? months : ["—"],
                datasets: [
                    {
                        label: "CO2 évité (kg)",
                        data: values.length ? values : [0],
                        backgroundColor: "#2d7a4f",
                        borderRadius: 6,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
                scales: {
                    y: { beginAtZero: true, ticks: { precision: 0 } },
                },
            },
        });
    }

    const materiauxCanvas = document.getElementById("materiauxChart");
    if (materiauxCanvas) {
        const labels = materiauxStats.value.map((m) => m.type_materiau);
        const counts = materiauxStats.value.map((m) => m.count);
        new Chart(materiauxCanvas, {
            type: "doughnut",
            data: {
                labels: labels.length ? labels : ["Aucune donnée"],
                datasets: [
                    {
                        data: counts.length ? counts : [1],
                        backgroundColor: counts.length
                            ? materiauxColors
                            : ["#e0e0e0"],
                        borderWidth: 0,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
            },
        });
    }
}

function loadChartJs() {
    return new Promise((resolve) => {
        if (typeof Chart !== "undefined") return resolve();
        const script = document.createElement("script");
        script.src =
            "https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.1/chart.umd.js";
        script.onload = resolve;
        document.head.appendChild(script);
    });
}

onMounted(() => {
    const id = userId();
    if (isNaN(id)) return;

    const headers = authHeaders();
    const role = sessionStorage.getItem("userRole") || "Prestataire";

    loadMateriauxRecherches();

    Promise.all([
        fetch(`${API_URL}/users/${id}/stats`, { headers })
            .then((r) => r.json())
            .then((d) => {
                stats.value = d;
                if (d.total_points)
                    sessionStorage.setItem("userScore", d.total_points);
            })
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/achats`, { headers })
            .then((r) => r.json())
            .then((d) => (achats.value = Array.isArray(d) ? d : []))
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/projets`, { headers })
            .then((r) => r.json())
            .then((d) => (projets.value = Array.isArray(d) ? d : []))
            .catch(() => {}),

        fetch(`${API_URL}/tips/role/${role}`)
            .then((r) => r.json())
            .then((d) => (tipDuJour.value = d))
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/notifications`, { headers })
            .then((r) => r.json())
            .then((d) => (latestNotification.value = d[0] || null))
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/eco-stats`, { headers })
            .then((r) => r.json())
            .then((d) => (ecoStats.value = { ...ecoStats.value, ...d }))
            .catch(() => {}),

        fetch(`${API_URL}/materiaux/stats`, { headers })
            .then((r) => r.json())
            .then((d) => (materiauxStats.value = Array.isArray(d) ? d : []))
            .catch(() => {}),

        loadAlertesPrioritaires(),

        fetch(`${API_URL}/users/${id}/abonnement`, { headers })
            .then((r) => r.json())
            .then((d) => (abonnement.value = { ...abonnement.value, ...d }))
            .catch(() => {}),

        fetch(`${API_URL}/users/${id}/factures`, { headers })
            .then((r) => r.json())
            .then((d) => (factures.value = Array.isArray(d) ? d : []))
            .catch(() => {}),

        fetch(`${API_URL}/user/planning/${id}`, { headers })
            .then((r) => r.json())
            .then((d) => {
                if (!d || Array.isArray(d)) return;
                const mapE = (arr, kind, key) =>
                    (arr || [])
                        .map((i) => {
                            const p = new Date(i[key]);
                            return isNaN(p)
                                ? null
                                : {
                                      id: `${kind}-${i.id}`,
                                      kind,
                                      title: i.titre,
                                      date: getLocalISODate(p),
                                      timeLabel: p.toLocaleTimeString("fr-FR", {
                                          hour: "2-digit",
                                          minute: "2-digit",
                                      }),
                                  };
                        })
                        .filter(Boolean);
                calendarEntries.value = [
                    ...mapE(d.formations, "formation", "date_debut"),
                    ...mapE(d.evenements, "event", "date_evenement"),
                ].sort((a, b) => a.date.localeCompare(b.date));
            })
            .catch(() => {}),
    ]).then(async () => {
        await loadChartJs();
        await nextTick();
        renderCharts();
    });
});
</script>

<style scoped>
.public-dashboard {
    padding: 20px;
    background: #f7f9f7;
    min-height: 100vh;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 2rem;
}

.header-actions-group {
    display: flex;
    align-items: center;
    gap: 12px;
}

.premium-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: #fff3cd;
    color: #856404;
    padding: 8px 14px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: 800;
}

.hero-title1 {
    font-size: 2.2rem;
    font-weight: 900;
    color: #122018;
    margin: 0;
}

.sidebar__category2 {
    margin: 0;
    color: #a0ada7;
    font-family: "Space Mono", monospace;
    font-size: 0.65rem;
    text-transform: uppercase;
}

.classic-text {
    font-size: 0.95rem;
    color: #63746a;
    margin-top: 8px;
    line-height: 1.5;
}

.btn-main-action,
.btn-secondary,
.btn-text-green,
.btn-view {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    text-decoration: none;
    font-weight: bold;
    border-radius: 10px;
    transition: 0.2s;
}

.btn-main-action {
    background: #2d7a4f;
    color: white;
    padding: 10px 20px;
    border: none;
}

.btn-main-action:hover:not(:disabled) {
    background: #1b4d31;
}

.btn-main-action:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-secondary {
    padding: 8px 16px;
    border: 1px solid #ddd;
    background: white;
    color: #333;
}

.btn-secondary:hover {
    background: #f0f0f0;
}

.btn-view {
    background: transparent;
    border: 1px solid #ddd;
    padding: 6px 12px;
    border-radius: 6px;
    color: #333;
}

.btn-view:hover {
    background: #f0f0f0;
}

.btn-text-green {
    background: #eaf4ed;
    color: #2d7a4f;
    padding: 12px 20px;
    width: 100%;
    border: none;
}

.btn-text-green:hover {
    background: #d1e7dd;
    transform: translateY(-2px);
}

.stats-grid {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: 1.4fr 0.7fr 0.7fr 0.9fr;
    margin-bottom: 2rem;
}

@media (max-width: 1100px) {
    .stats-grid {
        grid-template-columns: 1fr 1fr;
    }
}

@media (max-width: 600px) {
    .stats-grid {
        grid-template-columns: 1fr;
    }
}

.card,
.section-container {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
    margin-bottom: 2rem;
}

.card--score {
    background: #2d7a4f;
    color: white;
    border: none;
}

.tag-score {
    background: rgba(255, 255, 255, 0.2);
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: 800;
    display: inline-block;
    margin-bottom: 1rem;
}

.score-value {
    font-size: 2.8rem;
    font-weight: 900;
    line-height: 1;
    margin-bottom: 0.5rem;
}

.score-value span {
    font-size: 1.2rem;
    font-weight: normal;
    opacity: 0.8;
}

.score-level {
    font-size: 1.1rem;
    margin-bottom: 1.5rem;
    opacity: 0.9;
}

.score-footer {
    display: flex;
    justify-content: space-between;
    border-top: 1px solid rgba(255, 255, 255, 0.2);
    padding-top: 1rem;
    margin-top: 1rem;
}

.mini-stat strong {
    font-size: 1.1rem;
    display: block;
}

.card-num {
    font-size: 2.5rem;
    font-weight: bold;
    color: #1a1a1a;
    margin-bottom: 0.5rem;
    line-height: 1;
}

.text-dm {
    color: #666;
    margin-bottom: 1rem;
}

.badge {
    padding: 4px 10px;
    border-radius: 20px;
    font-size: 0.8rem;
    font-weight: bold;
    display: inline-flex;
    align-items: center;
    gap: 4px;
}

.badge--green {
    background: #eaf4ed;
    color: #2d7a4f;
}

.badge--orange {
    background: #fff3cd;
    color: #e8a030;
}

.status-logistique {
    background: #f0f0f0;
    color: #555;
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
    text-transform: uppercase;
}

.status-valid {
    background: #e9f5ed;
    color: #1e5636;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}

.status-neutral {
    background: #f5f5f5;
    color: #666;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
}

.data-table {
    width: 100%;
    border-collapse: collapse;
}

.data-table th {
    padding: 12px;
    text-align: left;
    border-bottom: 2px solid #eee;
    color: #666;
    font-weight: 600;
}

.data-table td {
    padding: 16px 12px;
    border-bottom: 1px solid #eee;
}

.table-subtext {
    color: #666;
    font-size: 0.8rem;
    display: block;
    margin-top: 4px;
}

.text-truncate {
    display: inline-block;
    max-width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: #6d7b72;
}

.material-tag {
    background: #e9f5ed;
    color: #1e5636;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 700;
    display: inline-block;
}

.stats-text {
    font-size: 0.85rem;
    color: #555;
}

.actions-cell {
    display: flex;
    gap: 8px;
    align-items: center;
}

.text-right {
    text-align: right !important;
}

.text-center {
    text-align: center;
}

.text-grey {
    color: #999;
}

.py-4 {
    padding-top: 1rem;
    padding-bottom: 1rem;
}

.eco-metrics-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: 1.5rem;
}

@media (max-width: 900px) {
    .eco-metrics-grid {
        grid-template-columns: repeat(2, 1fr);
    }
}

.eco-metric-card {
    background: #fafdfb;
    border: 1px solid #eee;
    border-radius: 12px;
    padding: 1.2rem;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.eco-metric-label {
    font-size: 0.78rem;
    color: #888;
    text-transform: uppercase;
    font-weight: 700;
    letter-spacing: 0.5px;
}

.eco-metric-value {
    font-size: 1.6rem;
    font-weight: 800;
    color: #1a1a1a;
}

.eco-metric-trend {
    font-size: 0.78rem;
    display: flex;
    align-items: center;
    gap: 4px;
    font-weight: 600;
}

.trend-up {
    color: #2d7a4f;
}

.trend-down {
    color: #d32f2f;
}

.chart-wrapper {
    background: #fafdfb;
    border: 1px solid #eee;
    border-radius: 12px;
    padding: 1.2rem;
}

.materials-grid {
    display: grid;
    grid-template-columns: 1fr 1.4fr;
    gap: 1.5rem;
    align-items: start;
}

@media (max-width: 900px) {
    .materials-grid {
        grid-template-columns: 1fr;
    }
}

.materials-table-wrapper {
    background: #fafdfb;
    border: 1px solid #eee;
    border-radius: 12px;
    padding: 1.2rem;
}

.chart-legend {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    margin-top: 1rem;
    justify-content: center;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.8rem;
    color: #555;
}

.legend-dot {
    width: 10px;
    height: 10px;
    border-radius: 3px;
    display: inline-block;
}

.materiaux-form {
    display: flex;
    gap: 12px;
    align-items: center;
}

.materiaux-input {
    flex: 1;
    padding: 0.9rem 1rem;
    border: 1px solid #ddd;
    border-radius: 10px;
    font-family: inherit;
    font-size: 0.95rem;
}

.materiaux-input:focus {
    outline: none;
    border-color: #2d7a4f;
}

.materiaux-saved-msg {
    margin-top: 0.8rem;
    color: #2d7a4f;
    font-weight: 700;
    font-size: 0.9rem;
}

.materiaux-keywords-preview {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 1rem;
}

.planning-week {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 1rem;
}

.planning-day {
    background: #fbfdfb;
    border: 1px solid #e0e0e0;
    border-radius: 12px;
    padding: 10px;
    min-height: 100px;
}

.planning-day.is-today {
    border: 2px solid #2d7a4f;
    background: #eaf4ed;
}

.day-label {
    text-align: center;
    color: #666;
    font-size: 0.85rem;
    margin-bottom: 8px;
    display: flex;
    flex-direction: column;
}

.day-label strong {
    color: #1a1a1a;
    font-size: 1.1rem;
    margin-top: 4px;
}

.entry {
    padding: 6px;
    border-radius: 6px;
    font-size: 0.75rem;
    margin-bottom: 4px;
    display: flex;
    flex-direction: column;
}

.entry--formation {
    background: #d1e7dd;
    color: #0f5132;
}

.entry--event {
    background: #fff3cd;
    color: #856404;
}

.end-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
}

.mt-2 {
    margin-top: 0.5rem;
}

.mt-4 {
    margin-top: 1rem;
}

.mb-2 {
    margin-bottom: 0.5rem;
}

.state-card {
    border: 1px dashed #cfe0d4;
    border-radius: 14px;
    padding: 26px;
    text-align: center;
    color: #666;
}
</style>
