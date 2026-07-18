"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  Camera,
  Scan,
  CheckCircle2,
  XCircle,
  RefreshCw,
  AlertTriangle,
  Info,
  ArrowRight,
  Shield,
} from "lucide-react";
import { toast } from "sonner";
import { FONT_CLASSES, Logo, Footer } from "../components/site-ui";

// ─── Types ────────────────────────────────────────────────────────────────────
type ScanState = "idle" | "camera" | "scanning" | "success" | "error" | "no_permission";

// ─── Helpers ──────────────────────────────────────────────────────────────────
function dataURLtoBlob(dataURL: string): Blob {
  const [header, data] = dataURL.split(",");
  const mime = header.match(/:(.*?);/)![1];
  const binary = atob(data);
  const arr = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) arr[i] = binary.charCodeAt(i);
  return new Blob([arr], { type: mime });
}

// ─── Scanner Overlay ─────────────────────────────────────────────────────────
function ScanOverlay({ scanning }: { scanning: boolean }) {
  return (
    <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
      {/* Dim corners */}
      <div className="absolute inset-0 bg-[#1F1B16]/50" style={{ clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%, 0 30%, 15% 30%, 15% 70%, 85% 70%, 85% 30%, 0 30%)" }} />

      {/* Face oval guide */}
      <div
        className={`relative w-56 h-72 rounded-[50%] border-2 transition-colors duration-500 ${
          scanning ? "border-[#B8863C]" : "border-white/80"
        }`}
      >
        {/* Corner accents */}
        {[
          "top-0 left-0 border-t-[3px] border-l-[3px] rounded-tl-[50%]",
          "top-0 right-0 border-t-[3px] border-r-[3px] rounded-tr-[50%]",
          "bottom-0 left-0 border-b-[3px] border-l-[3px] rounded-bl-[50%]",
          "bottom-0 right-0 border-b-[3px] border-r-[3px] rounded-br-[50%]",
        ].map((cls, i) => (
          <div
            key={i}
            className={`absolute w-8 h-8 border-[#B8863C] ${cls}`}
          />
        ))}

        {/* Scan line animation */}
        {scanning && (
          <div
            className="absolute left-2 right-2 h-0.5 bg-gradient-to-r from-transparent via-[#B8863C] to-transparent"
            style={{
              animation: "scan-line 1.8s ease-in-out infinite",
              top: "50%",
            }}
          />
        )}
      </div>

      {/* Instruction label */}
      <div className="absolute bottom-12 left-1/2 -translate-x-1/2 w-max max-w-xs text-center">
        <span className="bg-[#1F1B16]/80 backdrop-blur-sm text-[#F6F3EC] text-xs px-4 py-2 rounded-full font-[family-name:var(--font-mono)]">
          {scanning ? "Sedang memindai wajah…" : "Posisikan wajah di dalam oval"}
        </span>
      </div>
    </div>
  );
}

// ─── Main Component ───────────────────────────────────────────────────────────
export default function FaceScanPage() {
  const router = useRouter();
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);

  const [scanState, setScanState] = useState<ScanState>("idle");
  const [capturedImage, setCapturedImage] = useState<string | null>(null);
  const [matchedPhotos, setMatchedPhotos] = useState<number>(0);
  const [eventSlug, setEventSlug] = useState<string>("");

  // Read event from query param on mount
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    setEventSlug(params.get("event") ?? "");
  }, []);

  // ─── Camera control ────────────────────────────────────────────────────────
  const startCamera = useCallback(async () => {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: "user", width: { ideal: 1280 }, height: { ideal: 720 } },
        audio: false,
      });
      streamRef.current = stream;
      if (videoRef.current) {
        videoRef.current.srcObject = stream;
        videoRef.current.play();
      }
      setScanState("camera");
    } catch (err: unknown) {
      const e = err as Error;
      if (e.name === "NotAllowedError" || e.name === "PermissionDeniedError") {
        setScanState("no_permission");
      } else {
        toast.error("Kamera tidak dapat diakses. Pastikan perangkat Anda memiliki kamera.");
        setScanState("error");
      }
    }
  }, []);

  const stopCamera = useCallback(() => {
    streamRef.current?.getTracks().forEach((t) => t.stop());
    streamRef.current = null;
  }, []);

  // Cleanup on unmount
  useEffect(() => () => stopCamera(), [stopCamera]);

  // ─── Capture & (mock) scan ─────────────────────────────────────────────────
  const handleScan = useCallback(async () => {
    if (!videoRef.current || !canvasRef.current) return;
    setScanState("scanning");

    const video = videoRef.current;
    const canvas = canvasRef.current;
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    canvas.getContext("2d")?.drawImage(video, 0, 0);
    const dataURL = canvas.toDataURL("image/jpeg", 0.85);
    setCapturedImage(dataURL);
    stopCamera();

    // ── TODO (Fase 2): replace this mock with a real face-recognition API call ──
    // const blob = dataURLtoBlob(dataURL);
    // const form = new FormData();
    // form.append("face", blob, "face.jpg");
    // if (eventSlug) form.append("event_slug", eventSlug);
    // const token = localStorage.getItem("client_token");
    // const res = await fetch(`${API_URL}/client/face-search`, {
    //   method: "POST",
    //   headers: token ? { Authorization: `Bearer ${token}` } : {},
    //   body: form,
    // });
    // const data = await res.json();
    // ── End TODO ──

    // Mock: simulate 2 s processing then return 3 matched photos
    await new Promise((r) => setTimeout(r, 2000));
    setMatchedPhotos(3);
    setScanState("success");
    // Suppress linting for unused variable in TODO section
    void dataURLtoBlob;
  }, [stopCamera]);

  // ─── Reset ─────────────────────────────────────────────────────────────────
  const handleReset = () => {
    setCapturedImage(null);
    setMatchedPhotos(0);
    setScanState("idle");
  };

  // ─── Render ────────────────────────────────────────────────────────────────
  return (
    <>
      {/* Scan line keyframe */}
      <style>{`
        @keyframes scan-line {
          0%, 100% { top: 10%; opacity: 0; }
          10% { opacity: 1; }
          50% { top: 90%; opacity: 1; }
          90% { opacity: 1; }
        }
        @keyframes pulse-ring {
          0% { transform: scale(1); opacity: 0.6; }
          100% { transform: scale(1.4); opacity: 0; }
        }
      `}</style>

      <main className={`${FONT_CLASSES} min-h-screen flex flex-col bg-[#F6F3EC]`}>
        {/* Header */}
        <header className="sticky top-0 z-40 bg-[#F6F3EC]/95 backdrop-blur-sm border-b border-[#1F1B16]/8 px-4 py-3">
          <div className="max-w-6xl mx-auto flex items-center justify-between gap-4">
            <Logo small />
            <nav className="flex items-center gap-4">
              <Link href="/events" className="text-sm font-medium text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors">
                Event
              </Link>
              <Link href="/login" className="text-sm font-medium text-[#4A443C]/70 hover:text-[#1F1B16] transition-colors">
                Masuk
              </Link>
            </nav>
          </div>
        </header>

        <div className="flex-1 flex flex-col items-center justify-center px-4 py-12">
          <div className="w-full max-w-lg">
            {/* Page title */}
            <div className="text-center mb-8">
              <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#B8863C]/10 border border-[#B8863C]/20 text-[#B8863C] text-xs font-[family-name:var(--font-mono)] mb-4">
                <Shield size={12} />
                Fase 2 · Face Recognition
              </div>
              <h1 className="text-3xl font-[family-name:var(--font-display)] text-[#1F1B16] font-medium mb-2">
                Temukan Foto Anda
              </h1>
              <p className="text-[#4A443C]/70 text-sm leading-relaxed max-w-sm mx-auto">
                Pindai wajah Anda dan kami akan menemukan semua foto Anda dari{" "}
                {eventSlug ? (
                  <strong className="text-[#1F1B16]">event ini</strong>
                ) : (
                  "galeri event"
                )}{" "}
                secara otomatis.
              </p>
            </div>

            {/* ── IDLE ── */}
            {scanState === "idle" && (
              <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-lg overflow-hidden">
                <div className="p-8 text-center">
                  {/* Icon */}
                  <div className="relative w-24 h-24 mx-auto mb-6">
                    <div className="w-24 h-24 rounded-full bg-[#1F1B16]/5 flex items-center justify-center">
                      <Scan size={40} className="text-[#1F1B16]/40" />
                    </div>
                    {/* Pulse rings */}
                    {[1, 2].map((i) => (
                      <div
                        key={i}
                        className="absolute inset-0 rounded-full border border-[#B8863C]/30"
                        style={{ animation: `pulse-ring 2s ease-out ${i * 0.7}s infinite` }}
                      />
                    ))}
                  </div>

                  <h2 className="text-lg font-[family-name:var(--font-display)] text-[#1F1B16] mb-2">
                    Siap untuk Scan?
                  </h2>
                  <p className="text-sm text-[#4A443C]/70 mb-6 leading-relaxed">
                    Klik tombol di bawah untuk membuka kamera dan memindai wajah Anda. Pastikan pencahayaan cukup.
                  </p>

                  <button
                    id="btn-start-camera"
                    onClick={startCamera}
                    className="w-full py-4 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium flex items-center justify-center gap-2.5 hover:bg-[#B8863C] transition-all duration-300"
                  >
                    <Camera size={18} />
                    Buka Kamera
                  </button>
                </div>

                {/* Tips */}
                <div className="border-t border-[#1F1B16]/8 px-8 py-5 bg-[#F6F3EC]/50">
                  <p className="text-xs font-medium text-[#4A443C] mb-3 font-[family-name:var(--font-mono)]">
                    Tips agar hasil optimal:
                  </p>
                  <ul className="space-y-1.5">
                    {[
                      "Pastikan wajah terlihat jelas dan tidak tertutup",
                      "Gunakan tempat dengan pencahayaan yang baik",
                      "Hadapkan wajah langsung ke kamera",
                      "Lepaskan kacamata hitam atau masker",
                    ].map((tip) => (
                      <li key={tip} className="flex items-start gap-2 text-xs text-[#4A443C]/70">
                        <Info size={11} className="text-[#B8863C] shrink-0 mt-0.5" />
                        {tip}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            )}

            {/* ── CAMERA ── */}
            {(scanState === "camera" || scanState === "scanning") && (
              <div className="bg-[#1F1B16] rounded-3xl overflow-hidden shadow-2xl">
                {/* Video feed */}
                <div className="relative aspect-[3/4] bg-black">
                  <video
                    ref={videoRef}
                    className="w-full h-full object-cover"
                    playsInline
                    muted
                  />
                  <ScanOverlay scanning={scanState === "scanning"} />
                </div>

                {/* Controls */}
                <div className="p-6 flex items-center justify-between">
                  <button
                    id="btn-cancel-camera"
                    onClick={() => { stopCamera(); handleReset(); }}
                    className="flex items-center gap-2 text-[#F6F3EC]/60 hover:text-[#F6F3EC] transition-colors text-sm"
                  >
                    <XCircle size={16} />
                    Batal
                  </button>

                  {/* Shutter button */}
                  <button
                    id="btn-capture-scan"
                    onClick={handleScan}
                    disabled={scanState === "scanning"}
                    className="w-16 h-16 rounded-full bg-white border-4 border-[#B8863C] flex items-center justify-center hover:bg-[#F6F3EC] disabled:opacity-50 disabled:cursor-not-allowed transition-all active:scale-95 shadow-lg"
                  >
                    {scanState === "scanning" ? (
                      <div className="w-6 h-6 border-2 border-[#B8863C]/30 border-t-[#B8863C] rounded-full animate-spin" />
                    ) : (
                      <Scan size={22} className="text-[#B8863C]" />
                    )}
                  </button>

                  <div className="w-20" /> {/* Spacer */}
                </div>
              </div>
            )}

            {/* ── SUCCESS ── */}
            {scanState === "success" && (
              <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-lg overflow-hidden">
                {/* Captured preview */}
                {capturedImage && (
                  <div className="relative">
                    <img
                      src={capturedImage}
                      alt="Captured face"
                      className="w-full aspect-video object-cover"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-[#1F1B16]/80 to-transparent flex items-end p-6">
                      <div className="flex items-center gap-2 text-white">
                        <CheckCircle2 size={20} className="text-green-400" />
                        <span className="font-medium">Wajah berhasil dipindai</span>
                      </div>
                    </div>
                  </div>
                )}

                <div className="p-8 text-center">
                  {/* Result stat */}
                  <div className="w-16 h-16 rounded-2xl bg-green-50 border border-green-100 flex items-center justify-center mx-auto mb-4">
                    <CheckCircle2 size={28} className="text-green-500" />
                  </div>
                  <h2 className="text-xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-1">
                    {matchedPhotos} Foto Ditemukan!
                  </h2>
                  <p className="text-sm text-[#4A443C]/70 mb-6">
                    Kami menemukan <strong>{matchedPhotos}</strong> foto yang menampilkan wajah Anda. Lihat hasilnya di galeri.
                  </p>

                  {/* Note: fase 2 placeholder */}
                  <div className="bg-[#B8863C]/8 border border-[#B8863C]/20 rounded-xl p-4 mb-6 text-left">
                    <div className="flex items-start gap-2">
                      <AlertTriangle size={14} className="text-[#B8863C] shrink-0 mt-0.5" />
                      <p className="text-xs text-[#4A443C] leading-relaxed">
                        <span className="font-semibold text-[#1F1B16]">Fase 2 — Demo Mode:</span> Hasil di atas menggunakan data simulasi. Integrasi face recognition sebenarnya akan aktif pada peluncuran Fase 2.
                      </p>
                    </div>
                  </div>

                  <div className="flex flex-col sm:flex-row gap-3">
                    <button
                      id="btn-view-results"
                      onClick={() =>
                        eventSlug
                          ? router.push(`/events/${eventSlug}?face_matched=true`)
                          : router.push("/events")
                      }
                      className="flex-1 py-3.5 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm flex items-center justify-center gap-2 hover:bg-[#B8863C] transition-all duration-300"
                    >
                      Lihat Hasil <ArrowRight size={15} />
                    </button>
                    <button
                      id="btn-scan-again"
                      onClick={handleReset}
                      className="flex-1 py-3.5 rounded-2xl border border-[#1F1B16]/12 text-[#1F1B16] font-medium text-sm flex items-center justify-center gap-2 hover:border-[#B8863C] hover:text-[#B8863C] transition-all"
                    >
                      <RefreshCw size={15} />
                      Scan Ulang
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* ── NO PERMISSION ── */}
            {scanState === "no_permission" && (
              <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-lg p-8 text-center">
                <div className="w-16 h-16 rounded-2xl bg-amber-50 border border-amber-100 flex items-center justify-center mx-auto mb-4">
                  <AlertTriangle size={28} className="text-amber-500" />
                </div>
                <h2 className="text-xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-2">
                  Izin Kamera Ditolak
                </h2>
                <p className="text-sm text-[#4A443C]/70 mb-6 leading-relaxed">
                  Browser Anda memblokir akses kamera. Untuk menggunakan fitur ini, izinkan akses kamera di pengaturan browser Anda lalu muat ulang halaman.
                </p>
                <button
                  id="btn-reload-page"
                  onClick={() => window.location.reload()}
                  className="px-6 py-3 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] transition-all duration-300"
                >
                  Muat Ulang Halaman
                </button>
              </div>
            )}

            {/* ── ERROR ── */}
            {scanState === "error" && (
              <div className="bg-white rounded-3xl border border-[#1F1B16]/8 shadow-lg p-8 text-center">
                <div className="w-16 h-16 rounded-2xl bg-red-50 border border-red-100 flex items-center justify-center mx-auto mb-4">
                  <XCircle size={28} className="text-red-500" />
                </div>
                <h2 className="text-xl font-[family-name:var(--font-display)] text-[#1F1B16] mb-2">
                  Terjadi Kesalahan
                </h2>
                <p className="text-sm text-[#4A443C]/70 mb-6">
                  Kamera tidak dapat diakses. Pastikan perangkat Anda memiliki kamera yang berfungsi.
                </p>
                <button
                  id="btn-retry-camera"
                  onClick={handleReset}
                  className="px-6 py-3 rounded-2xl bg-[#1F1B16] text-[#F6F3EC] font-medium text-sm hover:bg-[#B8863C] transition-all duration-300"
                >
                  Coba Lagi
                </button>
              </div>
            )}

            {/* Hidden canvas for capture */}
            <canvas ref={canvasRef} className="hidden" />

            {/* Privacy notice */}
            <p className="text-center text-xs text-[#4A443C]/50 mt-6 font-[family-name:var(--font-mono)] leading-relaxed">
              <Shield size={10} className="inline mr-1" />
              Foto wajah Anda tidak disimpan. Hanya digunakan untuk pencarian foto secara real-time.
            </p>
          </div>
        </div>

        <Footer />
      </main>
    </>
  );
}
