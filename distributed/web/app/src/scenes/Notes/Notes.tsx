import React from 'react';
import { connect } from 'react-redux';
import Call, { Note } from '../../api';
import NotesList from './components/NotesList';
import NotesEditor from './components/NotesEditor';
import PageLayout from '../../components/PageLayout';
import './Notes.scss';
import { setNotes } from '../../store/notes';

interface Props {
  history: any;
  match: any;

  notes: Note[];
  setNotes: (notes: Note[]) => void;
}

class NotesScene extends React.Component<Props> {
  _mounted = false;

  componentDidMount() {
    this._mounted = true;
    const { notes, match, history, setNotes } = this.props;

    // Set the default note when navigating to /notes
    if(!match.params.id) {
      const id = notes.length === 0 ? 'new' : notes[0].id;
      history.push('/distributed/notes/' + id);
      return
    }

    // Fetch the notes from the API if we have none
    if(notes.length > 0) return;
    Call('listNotes').catch(console.warn).then(res => {
      if(!this._mounted || !res) return;
      
      const notes = (res.data.notes || []).map((n: any) => {
        return new Note(n);
      });

      setNotes(notes);
    })
  }

  componentWillUnmount() {
    this._mounted = false;
  }

  render():JSX.Element {
    const { notes } = this.props;
    const activeNoteID = this.props.match.params.id;
    const autoFocus = this.props.match.params.options === 'autoFocus';

    return(
      <PageLayout className='NotesScene'>
        <div className='notes-upper'>
          <h1>Notes</h1>
          <p>There are {notes.length} notes</p>
        </div>

        <div className='notes-lower'>
          <NotesList
            notes={notes}
            activeNoteID={activeNoteID}
            onNoteClicked={this.onNoteClicked.bind(this)} />
            
          <NotesEditor
            key={activeNoteID}
            noteID={activeNoteID}
            autoFocus={autoFocus}
            history={this.props.history} />
        </div>
      </PageLayout>
    );
  }

  onNoteClicked(id: string) {
    this.props.history.push('/distributed/notes/' + id)
  }
}

function mapDispatchToProps(dispatch: Function):any {
  return {
    setNotes: (notes: Note[]) => dispatch(setNotes(notes)),
  };
}

function mapStateToProps(state: any):any {
  return {
    notes: Object.values(state.notes.notes),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(NotesScene);