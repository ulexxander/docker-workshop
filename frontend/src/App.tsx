import { createEffect, createEvent, createStore, forward } from "effector";
import { useStore } from "effector-react";
import { useEffect } from "react";
import { apiRequest, Note, NoteCreateParams } from "./api";
import "./App.css";

const pageMounted = createEvent();

const $notes = createStore<Note[] | null>(null);
const $notesError = createStore<Error | null>(null);

const notesLoadFx = createEffect<void, Note[]>((params) =>
  apiRequest("/notes/all", { method: "GET" })
);

const noteCreateFx = createEffect<NoteCreateParams, Note>((params) =>
  apiRequest("/notes/create", {
    method: "POST",
    body: JSON.stringify(params),
  })
);

$notes.on(notesLoadFx.doneData, (_, notes) => notes);
$notesError.on(notesLoadFx.failData, (_, err) => err).reset(notesLoadFx);

forward({
  from: pageMounted,
  to: notesLoadFx,
});

forward({
  from: noteCreateFx.doneData,
  to: notesLoadFx,
});

const NotesList: React.FC = () => {
  const notes = useStore($notes);
  const notesError = useStore($notesError);

  if (notesError) {
    return <p>Error: {notesError.message}</p>;
  }
  if (!notes) {
    return <p>Loading...</p>;
  }
  if (!notes.length) {
    return <p>No notes yet...</p>;
  }

  return (
    <ul>
      {notes.map((note) => (
        <li key={note.ID}>
          <p>{note.Text}</p>
          <p>Created: {new Date(note.CreatedAt).toString()}</p>
        </li>
      ))}
    </ul>
  );
};

const Notes: React.FC = () => {
  return (
    <div>
      <h2>My Notes</h2>
      <NotesList />
    </div>
  );
};

export const App: React.FC = () => {
  useEffect(() => {
    pageMounted();
  }, []);

  return (
    <div className="App">
      <header>
        <img
          src="https://www.docker.com/sites/default/files/d8/styles/large/public/2021-08/Moby-share.png?itok=Kc8zKIm4"
          alt="Docker is fun"
        />
        <h1>Docker Workshop</h1>
      </header>

      <Notes />
    </div>
  );
};
