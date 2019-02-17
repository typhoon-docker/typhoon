import React from 'react';

import Link from '/components/Link/';

import { main, img, title } from './Login.css';
import { saveLocation } from '/utils/connect';

import typhoonImg from '/public/docker-swarm-typhoon.196x196.png';

const Login = () => (
  <div className={main}>
    <img src={typhoonImg} className={img} alt="Icône typhoon" />
    <h1 className={title}>Déployez vos sites simplement avec Typhoon</h1>
    <p>
      {
        "Typhoon s'occupe tout seul d'aller récupérer votre code automatiquement et le rend disponible pour vous sur internet."
      }
    </p>
    <p>
      <Link
        color="primary"
        bold
        href={`${process.env.BACKEND_URL}/login/viarezo`}
        title="Connecte toi"
        onClick={saveLocation}
      >
        Se connecter
      </Link>{' '}
      avec mon compte{' '}
      <Link color="viarezo" href="https://viarezo.fr/" title="Site web de ViaRezo" newTab>
        ViaRezo
      </Link>
      .
    </p>
  </div>
);

export default Login;
