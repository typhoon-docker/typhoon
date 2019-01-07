import React from "react";

import Button from "/components/Button/";
import Link from "/components/Link/";

import { main, img, title } from "./Login.css";

import typhoonImg from "/public/docker-swarm-typhoon.196x196.png";

const Login = () => {
  return (
    <div className={main}>
      <img src={typhoonImg} className={img} alt="Icône typhoon" />
      <h1 className={title}>Déployez vos sites simplement avec Typhoon</h1>
      <p>
        Typhoon s'occupe tout seul d'aller récupérer votre code automatiquement et le rend disponible pour vous sur
        internet.
      </p>
      <p>
        <Link color="primary" bold href="https://viarezo.fr/" title="Connecte toi">
          Se connecter
        </Link>{" "}
        avec mon compte{" "}
        <Link color="viarezo" href="https://viarezo.fr/" title="Site web de ViaRezo" newTab>
          ViaRezo
        </Link>
        .
      </p>
    </div>
  );
};

export default Login;
