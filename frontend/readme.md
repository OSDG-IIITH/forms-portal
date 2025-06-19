## Forms Portal Frontend

> This directory contains the source code for the frontend of the Forms Portal.

The frontend is made with `sveltekit`, using `shadcn-svelte` and `tailwind`.


## Prerequisites

- [Node.js](https://nodejs.org/)


## Steps

Install dependencies:

```sh
npm i
```

Start the development server:

```sh
npm run dev
```

The app will be available at [http://localhost:5173](http://localhost:5173),


## Usage/Testing

Two demo pages are available as of right now:
- `/create/demo`: Form Builder Demo
- `/view/demo`: Form Viewer (Responder) Demo


## Implemented Functionality

For demo purposes, the form builder (`/create/demo`) lets you save the form to json, and the form viewer (`/view/demo`) loads up the form stored at `src/lib/form.json`
