import { createRouter, createWebHistory } from "vue-router";

import Register from "./pages/Register.vue";
import Login from "./pages/Login.vue";
import Profil from "./pages/Profil.vue";

import Home from "./pages/Home.vue";
import Annonces from "./pages/Annonces.vue";
import Forums from "./pages/Forums.vue";

import Formations from "./pages/Formations.vue";
import Evenements from "./pages/Evenements.vue";
import Projets from "./pages/Projets.vue";
import Conseils from "./pages/Conseils.vue";
import Panier from "./pages/PanierDetail.vue";
import ClaimObjet from "./pages/ClaimObjet.vue";
import Messages from "./pages/Messages.vue";
import AbonnementDM from "./pages/AbonnementDM.vue";
import Abonnement from "./pages/Abonnement.vue";
import Paiement from "./pages/Paiement.vue";

const routes = [
  { path: "/front", component: Home },
  { path: "/inscription", component: Register },
  { path: "/connexion", component: Login },
  { path: "/messages", component: Messages },
  { path: "/catalogue", component: Annonces },
  { path: "/forums", component: Forums },
  { path: "/formations", component: Formations },
  { path: "/evenements", component: Evenements },
  { path: "/projets", component: Projets },
  { path: "/conseils", component: Conseils },
  { path: "/panier", component: Panier },
  { path: "/paiement", component: Paiement },
  { path: "/abonnement-dm", component: AbonnementDM },
  { path: "/abonnement", component: Abonnement },
  { 
    path: '/:pathMatch(.*)*', 
    redirect: '/front' 
  },
  {
    path: '/annonce/:id',
    name: 'annonce-detail',
    component: () => import('./pages/AnnonceDetail.vue')
  },
  {
    path: '/user/:id',
    name: 'utilisateur-detail',
    component: () => import('./pages/UtilisateurDetail.vue')
  },
  {
    path: '/formations/:id',
    name: 'formation-detail',
    component: () => import('./pages/FormationDetail.vue')
  },
  {
    path: '/evenements/:id',
    name: 'evenement-detail',
    component: () => import('./pages/EvenementDetail.vue')
  },
  {
    path: '/projets/:id',
    name: 'projet-detail',
    component: () => import('./pages/ProjetDetail.vue')
  },
  {
    path: '/conseils/:id',
    name: 'conseil-detail',
    component: () => import('./pages/ConseilDetail.vue')
  },
  {
    path: "/deposer",
    name: "deposer-objet",
    component: () => import('./pages/DropObjet.vue'),
    meta: { requiresRole: "Salarie" }
  },
  {
    path: "/module/science/1",
    name: "module-science1",
    component: () => import('./pages/ReportDataMining.vue'),
    meta: { requiresRole: "Admin" }
  },
  {
    path: "/claim",
    name: "recuperer-objet",
    component: () => import('./pages/ClaimObjet.vue'),
    meta: { requiresRole: "Salarie" }
  },
  {
    path: "/profil",
    component: Profil,
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/DashboardHome.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/DashboardHome.vue');
          return import('./pages/profilpagesParticulier/DashboardHome.vue');
        }
      },
      {
        path: "planning",
        name: "mon-planning",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Planning.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Planning.vue');
          return import('./pages/profilpagesParticulier/Planning.vue');
        }
      },
      {
        path: "informations",
        name: "mes-informations",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Infos.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Infos.vue');
          return import('./pages/profilpagesParticulier/Infos.vue');
        }
      },
      {
        path: "notifications",
        name: "mes-notifications",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Notification.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Notification.vue');
          return import('./pages/profilpagesParticulier/Notification.vue');
        }
      },
      {
        path: "factures",
        name: "mes-factures",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Factures.vue');
          return import('./pages/profilpagesParticulier/Factures.vue');
        }
      },
      {
        path: "/profil/forum",
        name: "mes-forums",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Forums.vue');
          return import('./pages/profilpagesParticulier/Forums.vue');
        }
      },
      {
        path: "password",
        name: "modification-password",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/ChangePassword.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/ChangePassword.vue');
          return import('./pages/profilpagesParticulier/ChangePassword.vue');
        }
      },
      
      {
        path: "annonces",
        name: "mes-annonces",
        component: () => import('./pages/profilpagesParticulier/Annonces.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "depots",
        name: "mes-depots",
        component: () => import('./pages/profilpagesParticulier/Depots.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "projet-favoris",
        name: "projet-favoris",
        component: () => import('./pages/profilpagesParticulier/ProjetLikes.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "modifyAnnonce/:id",
        name: "modification-annonce",
        component: () => import('./pages/profilpagesParticulier/ModifyAnnonce.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "createAnnonce",
        name: "create-annonce",
        component: () => import('./pages/profilpagesParticulier/CreateAnnonce.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "reserveCasier/:id",
        name: "reserve-casier",
        component: () => import('./pages/profilpagesParticulier/ReserverCasier.vue'),
        meta: { requiresRole: "Particulier" }
      },

      {
        path: "favoris",
        name: "annonce-like",
        component: () => import("./pages/profilpagesPrestataire/AnnoncesLike.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "recuperations",
        name: "objet-recuperation",
        component: () => import("./pages/profilpagesPrestataire/ObjetRecuperation.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "projets",
        name: "mes-projet",
        component: () => import("./pages/profilpagesPrestataire/Projets.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "modifyProjet/:id",
        name: "modify-projet",
        component: () => import("./pages/profilpagesPrestataire/ModifyProjet.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "createProjet",
        name: "create-projet",
        component: () => import("./pages/profilpagesPrestataire/CreateProjet.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "tips",
        name: "mes-tips",
        component: () => import("./pages/profilpagesSalarie/Tips.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "evenements",
        name: "mes-evenements",
        component: () => import("./pages/profilpagesSalarie/Evenement.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "formations",
        name: "mes-formations",
        component: () => import("./pages/profilpagesSalarie/Formation.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "createTips",
        name: "create-tips",
        component: () => import("./pages/profilpagesSalarie/CreateTips.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "createEvenement",
        name: "create-evenement",
        component: () => import("./pages/profilpagesSalarie/CreateEvenement.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "createFormation",
        name: "create-formation",
        component: () => import("./pages/profilpagesSalarie/CreateFormation.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "modifyTips/:id",
        name: "modify-tips",
        component: () => import("./pages/profilpagesSalarie/ModifyTip.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "modifyEvenement/:id",
        name: "modify-evenement",
        component: () => import("./pages/profilpagesSalarie/ModifyEvenement.vue"),
        meta: { requiresRole: "Salarie" }
      },
      {
        path: "modifyFormation/:id",
        name: "modify-formation",
        component: () => import("./pages/profilpagesSalarie/ModifyFormation.vue"),
        meta: { requiresRole: "Salarie" }
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: "smooth" };
    }
    return { top: 0 };
  },
});

router.beforeEach(async (to, from, next) => {
  const token = sessionStorage.getItem("userToken");
  const id = sessionStorage.getItem("userId");
  const userRole = sessionStorage.getItem("userRole");

  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const requiredRole = to.matched.find((record) => record.meta.requiresRole)?.meta.requiresRole;
  const isProtected = requiresAuth || requiredRole;

  if (!isProtected) {
    return next();
  }

  if (!id || !token) {
    sessionStorage.clear(); 
    return next("/connexion");
  }

  try {
    const response = await fetch(`/go/check-session?id=${id}`, {
      headers: { Authorization: token },
    });
    const data = await response.json();

    if (!data.isValid) {
      sessionStorage.clear();
      return next("/connexion");
    }
  } catch (error) {
    console.error("Erreur de vérification de session:", error);
    return next("/front");
  }

  if (requiredRole && userRole !== requiredRole) {
    console.warn(`Accès refusé ! Rôle attendu : ${requiredRole}, Rôle actuel : ${userRole}`);
    
    if (userRole === "Prestataire" || userRole === "Salarie" || userRole === "Particulier") {
      return next("/profil"); 
    } else {
      return next("/front"); 
    }
  }

  next();
});

export default router;