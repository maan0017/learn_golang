import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { ChangeDetectorRef, Component, HostListener, signal } from '@angular/core';

@Component({
  selector: 'app-root',
  imports: [CommonModule],
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {
  protected readonly title = signal('frontend-angular');

  imageDetails = signal<{
    name: string;
    size: number;
    type: string;
    format: string;
    url: string | null;
  } | null>(null);

  uploading = signal<boolean>(false);

  constructor(private http: HttpClient) {}

  // @HostListener('window:keydown', ['$event'])
  // keyJobs(event: KeyboardEvent) {
  //   console.log('event: ', event.key);
  // }

  onFileSelection(event: Event) {
    console.log('event triggered');
    this.uploading.set(true); // Start loader

    const input = event.target as HTMLInputElement;
    const file = input?.files ? input.files[0] : null;
    if (!file) return;

    console.log('file: ', file);

    // Prepare FormData
    const form = new FormData();
    form.append('file', file);

    // Upload to backend
    this.http.post<{ url: string }>('http://localhost:8080/upload', form).subscribe({
      next: (res) => {
        console.log('Upload success:', res);

        // Once upload is done, update details with Cloudinary URL
        this.imageDetails.set({
          name: file.name,
          size: file.size,
          type: file.type,
          format: file.name.split('.').pop() || 'Unknown',
          url: res.url, // use server's Cloudinary secure URL
        });

        this.uploading.set(false); // Stop loader
      },
      error: (err) => {
        console.error('Upload failed:', err);
        this.uploading.set(false); // Stop loader even on error
      },
    });
  }
}
